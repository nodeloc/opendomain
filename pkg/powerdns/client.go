package powerdns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client PowerDNS API 客户端
type Client struct {
	BaseURL    string
	APIKey     string
	ServerID   string
	HTTPClient *http.Client
}

// NewClient 创建新的 PowerDNS 客户端
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:  baseURL,
		APIKey:   apiKey,
		ServerID: "localhost", // 默认服务器 ID
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Zone 表示一个 DNS Zone
type Zone struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Kind        string   `json:"kind"` // Master, Slave, Native
	DNSsec      bool     `json:"dnssec,omitempty"`
	Serial      uint32   `json:"serial,omitempty"`
	Nameservers []string `json:"nameservers,omitempty"`
	RRsets      []RRset  `json:"rrsets,omitempty"`
}

// RRset 表示资源记录集
type RRset struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	TTL        int       `json:"ttl"`
	ChangeType string    `json:"changetype"` // REPLACE, DELETE
	Records    []Record  `json:"records,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

// Record 表示单条 DNS 记录
type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

// Comment 表示记录注释
type Comment struct {
	Content    string `json:"content"`
	Account    string `json:"account,omitempty"`
	ModifiedAt int64  `json:"modified_at,omitempty"`
}

// CreateZone 创建新的 DNS Zone
func (c *Client) CreateZone(domain string, nameservers []string) error {
	zone := &Zone{
		Name:        ensureTrailingDot(domain),
		Kind:        "Master",
		Nameservers: nameservers,
	}

	url := fmt.Sprintf("%s/api/v1/servers/%s/zones", c.BaseURL, c.ServerID)
	body, err := json.Marshal(zone)
	if err != nil {
		return fmt.Errorf("failed to marshal zone: %w", err)
	}

	_, err = c.doRequest("POST", url, body)
	return err
}

// DeleteZone 删除 DNS Zone
func (c *Client) DeleteZone(domain string) error {
	zoneName := ensureTrailingDot(domain)
	url := fmt.Sprintf("%s/api/v1/servers/%s/zones/%s", c.BaseURL, c.ServerID, zoneName)
	_, err := c.doRequest("DELETE", url, nil)
	return err
}

// GetZone 获取 Zone 信息
func (c *Client) GetZone(domain string) (*Zone, error) {
	zoneName := ensureTrailingDot(domain)
	url := fmt.Sprintf("%s/api/v1/servers/%s/zones/%s", c.BaseURL, c.ServerID, zoneName)

	respBody, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var zone Zone
	if err := json.Unmarshal(respBody, &zone); err != nil {
		return nil, fmt.Errorf("failed to unmarshal zone: %w", err)
	}

	return &zone, nil
}

// RecordEntry 表示一条待同步的记录
type RecordEntry struct {
	Content  string
	Priority *int
}

// SetRecords 设置某个 name+type 的完整记录集（支持多条记录）
func (c *Client) SetRecords(domain, name, recordType string, entries []RecordEntry, ttl int) error {
	zoneName := ensureTrailingDot(domain)
	recordName := ensureTrailingDot(name)

	records := make([]Record, 0, len(entries))
	for _, e := range entries {
		var content string
		if recordType == "MX" && e.Priority != nil {
			content = fmt.Sprintf("%d %s", *e.Priority, ensureTrailingDot(e.Content))
		} else if recordType == "CNAME" || recordType == "NS" || recordType == "MX" {
			content = ensureTrailingDot(e.Content)
		} else {
			content = e.Content
		}
		records = append(records, Record{Content: content, Disabled: false})
	}

	rrset := RRset{
		Name:       recordName,
		Type:       recordType,
		TTL:        ttl,
		ChangeType: "REPLACE",
		Records:    records,
	}

	return c.patchRRset(zoneName, rrset)
}

// DeleteRRset 删除某个 name+type 的所有记录
func (c *Client) DeleteRRset(domain, name, recordType string) error {
	zoneName := ensureTrailingDot(domain)
	recordName := ensureTrailingDot(name)

	rrset := RRset{
		Name:       recordName,
		Type:       recordType,
		ChangeType: "DELETE",
	}

	return c.patchRRset(zoneName, rrset)
}

// patchRRsets 批量发送多个 RRset patch 请求
func (c *Client) patchRRsets(zoneName string, rrsets []RRset) error {
	patchData := map[string][]RRset{
		"rrsets": rrsets,
	}

	body, err := json.Marshal(patchData)
	if err != nil {
		return fmt.Errorf("failed to marshal rrsets: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/servers/%s/zones/%s", c.BaseURL, c.ServerID, zoneName)
	_, err = c.doRequest("PATCH", url, body)
	return err
}

// SetSubdomainDisabled 将某个子域名下的所有记录设为 disabled/enabled
// subdomain 是完整子域名（如 "test.loc.cc"），会匹配该子域名及其子记录
func (c *Client) SetSubdomainDisabled(rootDomain, subdomain string, disabled bool) error {
	zoneName := ensureTrailingDot(rootDomain)

	// 获取 zone 的所有 RRsets
	zone, err := c.GetZone(rootDomain)
	if err != nil {
		return fmt.Errorf("failed to get zone: %w", err)
	}

	fqdnSuffix := ensureTrailingDot(subdomain)
	var patches []RRset

	for _, rrset := range zone.RRsets {
		// 跳过 SOA/NS 等系统记录
		if rrset.Type == "SOA" || rrset.Type == "NS" {
			continue
		}
		// 匹配属于该子域名的记录（精确匹配或子域名记录）
		if rrset.Name == fqdnSuffix || strings.HasSuffix(rrset.Name, "."+fqdnSuffix) {
			newRecords := make([]Record, len(rrset.Records))
			for i, r := range rrset.Records {
				newRecords[i] = Record{
					Content:  r.Content,
					Disabled: disabled,
				}
			}
			patches = append(patches, RRset{
				Name:       rrset.Name,
				Type:       rrset.Type,
				TTL:        rrset.TTL,
				ChangeType: "REPLACE",
				Records:    newRecords,
			})
		}
	}

	if len(patches) == 0 {
		return nil // 没有需要更新的记录
	}

	return c.patchRRsets(zoneName, patches)
}

// patchRRset 发送 RRset patch 请求
func (c *Client) patchRRset(zoneName string, rrset RRset) error {
	patchData := map[string][]RRset{
		"rrsets": {rrset},
	}

	body, err := json.Marshal(patchData)
	if err != nil {
		return fmt.Errorf("failed to marshal rrset: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/servers/%s/zones/%s", c.BaseURL, c.ServerID, zoneName)
	_, err = c.doRequest("PATCH", url, body)
	return err
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(method, url string, body []byte) ([]byte, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ensureTrailingDot 确保域名以点结尾
func ensureTrailingDot(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[len(s)-1] != '.' {
		return s + "."
	}
	return s
}
