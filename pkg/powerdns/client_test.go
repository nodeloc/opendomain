package powerdns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCreateZoneEnablesAPIRectify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/servers/localhost/zones" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var zone Zone
		if err := json.NewDecoder(r.Body).Decode(&zone); err != nil {
			t.Fatalf("decode zone: %v", err)
		}
		if zone.Name != "child.loc.cc." || !zone.APIRectify {
			t.Fatalf("unexpected zone payload: %+v", zone)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.CreateZone("child.loc.cc", []string{"ns1.example.", "ns2.example."}); err != nil {
		t.Fatalf("CreateZone: %v", err)
	}
}

func TestEnsureDelegatedZoneCreatesChildBeforeParentDelegation(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Fatal("missing API key")
		}

		switch len(calls) {
		case 1:
			w.WriteHeader(http.StatusNotFound)
		case 2:
			var zone Zone
			if err := json.NewDecoder(r.Body).Decode(&zone); err != nil {
				t.Fatalf("decode zone: %v", err)
			}
			if zone.Name != "child.loc.cc." || !zone.APIRectify {
				t.Fatalf("unexpected child zone: %+v", zone)
			}
			w.WriteHeader(http.StatusCreated)
		case 3:
			var patch struct {
				RRsets []RRset `json:"rrsets"`
			}
			if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
				t.Fatalf("decode delegation: %v", err)
			}
			if len(patch.RRsets) != 2 || patch.RRsets[0].Type != "CNAME" || patch.RRsets[0].ChangeType != "DELETE" ||
				patch.RRsets[1].Type != "NS" || patch.RRsets[1].Name != "child.loc.cc." {
				t.Fatalf("unexpected delegation: %+v", patch.RRsets)
			}
			got := []string{patch.RRsets[1].Records[0].Content, patch.RRsets[1].Records[1].Content}
			want := []string{"ns1.example.", "ns2.example."}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("delegation nameservers = %v, want %v", got, want)
			}
			w.WriteHeader(http.StatusNoContent)
		case 4:
			_, _ = w.Write([]byte(`{"name":"loc.cc.","api_rectify":false}`))
		case 5:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected extra request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.EnsureDelegatedZone("child.loc.cc", "loc.cc", []string{"ns1.example", "ns2.example"}); err != nil {
		t.Fatalf("EnsureDelegatedZone: %v", err)
	}

	wantCalls := []string{
		"GET /api/v1/servers/localhost/zones/child.loc.cc.",
		"POST /api/v1/servers/localhost/zones",
		"PATCH /api/v1/servers/localhost/zones/loc.cc.",
		"GET /api/v1/servers/localhost/zones/loc.cc.",
		"PUT /api/v1/servers/localhost/zones/loc.cc./rectify",
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("calls = %v, want %v", calls, wantCalls)
	}
}

func TestRemoveDelegationDeletesDSBeforeNSInOnePatch(t *testing.T) {
	var patches [][]string
	var rectifies int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			var patch struct {
				RRsets []RRset `json:"rrsets"`
			}
			if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
				t.Fatalf("decode patch: %v", err)
			}
			var types []string
			for _, rr := range patch.RRsets {
				if rr.ChangeType != "DELETE" || rr.Name != "child.loc.cc." {
					t.Fatalf("unexpected rrset %+v", rr)
				}
				types = append(types, rr.Type)
			}
			patches = append(patches, types)
		case http.MethodPut:
			rectifies++
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.RemoveDelegation("loc.cc", "child.loc.cc"); err != nil {
		t.Fatalf("RemoveDelegation: %v", err)
	}
	if want := [][]string{{"DS", "NS"}}; !reflect.DeepEqual(patches, want) {
		t.Fatalf("patches = %v, want %v", patches, want)
	}
	if rectifies != 1 {
		t.Fatalf("zone without API-RECTIFY must be rectified once, got %d", rectifies)
	}
}

func TestRemoveDelegationSkipsRectifyWhenZoneHasAPIRectify(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if r.Method == http.MethodGet {
			if r.URL.Query().Get("rrsets") != "false" {
				t.Fatalf("zone info lookup must not transfer rrsets: %q", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"name":"loc.cc.","api_rectify":true}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.RemoveDelegation("loc.cc", "child.loc.cc"); err != nil {
		t.Fatalf("RemoveDelegation: %v", err)
	}
	wantCalls := []string{
		"PATCH /api/v1/servers/localhost/zones/loc.cc.",
		"GET /api/v1/servers/localhost/zones/loc.cc.",
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("calls = %v, want %v", calls, wantCalls)
	}
}

func TestWritesToSameZoneAreSerializedAcrossClients(t *testing.T) {
	var mu sync.Mutex
	inFlight := map[string]int{}
	maxInFlight := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`{"api_rectify":true}`))
			return
		}
		zone := zoneFromURL(r.URL.Path)
		mu.Lock()
		inFlight[zone]++
		if inFlight[zone] > maxInFlight[zone] {
			maxInFlight[zone] = inFlight[zone]
		}
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		inFlight[zone]--
		mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client := NewClient(server.URL, "test-key")
			if err := client.RemoveDelegation("serial-test.cc", fmt.Sprintf("c%d.serial-test.cc", i)); err != nil {
				t.Errorf("RemoveDelegation: %v", err)
			}
		}(i)
	}
	wg.Wait()
	if maxInFlight["serial-test.cc."] != 1 {
		t.Fatalf("expected writes to one zone to be serialized, max in flight = %d", maxInFlight["serial-test.cc."])
	}
}

func TestZoneFromURL(t *testing.T) {
	cases := map[string]string{
		"http://pdns:8081/api/v1/servers/localhost/zones/Loc.CC.":                  "loc.cc.",
		"http://pdns:8081/api/v1/servers/localhost/zones/loc.cc./rectify":          "loc.cc.",
		"http://pdns:8081/api/v1/servers/localhost/zones/loc.cc.?rrsets=false":     "loc.cc.",
		"http://pdns:8081/api/v1/servers/localhost/zones/child.loc.cc./cryptokeys": "child.loc.cc.",
		"http://pdns:8081/api/v1/servers/localhost/zones":                          "",
	}
	for in, want := range cases {
		if got := zoneFromURL(in); got != want {
			t.Errorf("zoneFromURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEnableDNSSECEnablesAPIRectifyCreatesKeyAndRectifies(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch len(calls) {
		case 1: // PUT zone metadata: api_rectify=true
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode zone update: %v", err)
			}
			if body["api_rectify"] != true {
				t.Fatalf("expected api_rectify=true, got %+v", body)
			}
			w.WriteHeader(http.StatusNoContent)
		case 2: // GET cryptokeys: none yet
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("[]"))
		case 3: // POST cryptokeys
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode key: %v", err)
			}
			if body["keytype"] != "csk" || body["active"] != true || body["published"] != true {
				t.Fatalf("unexpected key payload: %+v", body)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":1}`))
		case 4: // PUT rectify
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected extra request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.EnableDNSSEC("child.loc.cc"); err != nil {
		t.Fatalf("EnableDNSSEC: %v", err)
	}

	wantCalls := []string{
		"PUT /api/v1/servers/localhost/zones/child.loc.cc.",
		"GET /api/v1/servers/localhost/zones/child.loc.cc./cryptokeys",
		"POST /api/v1/servers/localhost/zones/child.loc.cc./cryptokeys",
		"PUT /api/v1/servers/localhost/zones/child.loc.cc./rectify",
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("calls = %v, want %v", calls, wantCalls)
	}
}

func TestEnableDNSSECReturnsRectifyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`[{"id":1,"keytype":"csk","active":true,"published":true}]`))
		case r.Method == http.MethodPut && r.URL.Path == "/api/v1/servers/localhost/zones/child.loc.cc./rectify":
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"error":"Zone is not DNSSEC enabled"}`))
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.EnableDNSSEC("child.loc.cc"); err == nil {
		t.Fatal("expected rectify failure to be returned")
	}
}

func TestPublishDSToParentZoneUsesShortTTLAndRectifiesParent(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch len(calls) {
		case 1: // GET child cryptokeys
			_, _ = w.Write([]byte(`[
				{"id":1,"keytype":"csk","active":true,"published":true,"ds":["12345 13 2 ABCDEF"]},
				{"id":2,"keytype":"csk","active":false,"published":true,"ds":["99999 13 2 OLD"]}
			]`))
		case 2: // PATCH parent DS
			var patch struct {
				RRsets []RRset `json:"rrsets"`
			}
			if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
				t.Fatalf("decode patch: %v", err)
			}
			rr := patch.RRsets[0]
			if rr.Type != "DS" || rr.Name != "child.loc.cc." || rr.TTL != dsRecordTTL || rr.ChangeType != "REPLACE" {
				t.Fatalf("unexpected DS rrset: %+v", rr)
			}
			if len(rr.Records) != 1 || rr.Records[0].Content != "12345 13 2 ABCDEF" {
				t.Fatalf("inactive key DS must not be published: %+v", rr.Records)
			}
			w.WriteHeader(http.StatusNoContent)
		case 3: // GET parent zone info
			_, _ = w.Write([]byte(`{"name":"loc.cc.","api_rectify":false}`))
		case 4: // PUT rectify parent
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected extra request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.PublishDSToParentZone("child.loc.cc", "loc.cc"); err != nil {
		t.Fatalf("PublishDSToParentZone: %v", err)
	}
	wantCalls := []string{
		"GET /api/v1/servers/localhost/zones/child.loc.cc./cryptokeys",
		"PATCH /api/v1/servers/localhost/zones/loc.cc.",
		"GET /api/v1/servers/localhost/zones/loc.cc.",
		"PUT /api/v1/servers/localhost/zones/loc.cc./rectify",
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("calls = %v, want %v", calls, wantCalls)
	}
}

func TestUnpublishDSFromParentZoneRectifiesParent(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.UnpublishDSFromParentZone("child.loc.cc", "loc.cc"); err != nil {
		t.Fatalf("UnpublishDSFromParentZone: %v", err)
	}
	wantCalls := []string{
		"PATCH /api/v1/servers/localhost/zones/loc.cc.",
		"GET /api/v1/servers/localhost/zones/loc.cc.",
		"PUT /api/v1/servers/localhost/zones/loc.cc./rectify",
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("calls = %v, want %v", calls, wantCalls)
	}
}

func TestGetRRsetFiltersByNameAndType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("rrset_name") != "child.loc.cc." {
			t.Fatalf("expected rrset_name filter, got %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"name":"loc.cc.","dnssec":true,"rrsets":[
			{"name":"child.loc.cc.","type":"NS","ttl":3600,"records":[{"content":"ns1.example.","disabled":false}]},
			{"name":"child.loc.cc.","type":"DS","ttl":300,"records":[{"content":"12345 13 2 ABCDEF","disabled":false}]}
		]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	records, err := client.GetRRset("loc.cc", "child.loc.cc", "DS")
	if err != nil {
		t.Fatalf("GetRRset: %v", err)
	}
	if len(records) != 1 || records[0].Content != "12345 13 2 ABCDEF" {
		t.Fatalf("unexpected DS records: %+v", records)
	}
	none, err := client.GetRRset("loc.cc", "child.loc.cc", "TXT")
	if err != nil || none != nil {
		t.Fatalf("expected no TXT rrset, got %v %v", none, err)
	}
}

func TestDisableDNSSECSendsDnssecFalse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v1/servers/localhost/zones/child.loc.cc." {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["dnssec"] != false {
			t.Fatalf("expected dnssec=false, got %+v", body)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.DisableDNSSEC("child.loc.cc"); err != nil {
		t.Fatalf("DisableDNSSEC: %v", err)
	}
}

func TestSetDelegationsBatchesAndRectifiesOnce(t *testing.T) {
	var patches [][]RRset
	var rectifies int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPatch:
			var patch struct {
				RRsets []RRset `json:"rrsets"`
			}
			if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
				t.Fatalf("decode patch: %v", err)
			}
			patches = append(patches, patch.RRsets)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/servers/localhost/zones/loc.cc.":
			_, _ = w.Write([]byte(`{"name":"loc.cc.","api_rectify":false}`))
		case r.Method == http.MethodPut && r.URL.Path == "/api/v1/servers/localhost/zones/loc.cc./rectify":
			rectifies++
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	// 60 delegations x 3 rrsets = 180; batches of 40 delegations -> 120 + 60
	var dels []Delegation
	for i := 0; i < 60; i++ {
		d := Delegation{Child: fmt.Sprintf("c%d.loc.cc", i), Nameservers: []string{"ns1.example", "ns2.example"}}
		if i%2 == 0 {
			d.DS = []string{"1 13 2 AA"}
		}
		dels = append(dels, d)
	}

	client := NewClient(server.URL, "test-key")
	if err := client.SetDelegations("loc.cc", dels); err != nil {
		t.Fatalf("SetDelegations: %v", err)
	}
	if len(patches) != 2 || len(patches[0]) != 120 || len(patches[1]) != 60 {
		t.Fatalf("unexpected batching: %d patches, sizes %v", len(patches), func() []int {
			var s []int
			for _, p := range patches {
				s = append(s, len(p))
			}
			return s
		}())
	}
	if rectifies != 1 {
		t.Fatalf("expected exactly one rectify, got %d", rectifies)
	}
	first := patches[0]
	if first[0].Type != "CNAME" || first[0].ChangeType != "DELETE" || first[0].Name != "c0.loc.cc." {
		t.Fatalf("legacy CNAME at the delegation point must be deleted first: %+v", first[0])
	}
	if first[1].Type != "NS" || first[1].Name != "c0.loc.cc." || first[1].Records[0].Content != "ns1.example." {
		t.Fatalf("unexpected NS rrset: %+v", first[1])
	}
	if first[2].Type != "DS" || first[2].ChangeType != "REPLACE" || first[2].TTL != dsRecordTTL {
		t.Fatalf("signed child must get DS REPLACE: %+v", first[2])
	}
	if first[5].Type != "DS" || first[5].ChangeType != "DELETE" {
		t.Fatalf("unsigned child must get DS DELETE: %+v", first[5])
	}
}

func TestSetDomainSuspendedDisablesChildRecordsAndLegacyParentRecordsOnly(t *testing.T) {
	var patches = map[string][]RRset{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/servers/localhost/zones/child.loc.cc.":
			_, _ = w.Write([]byte(`{"name":"child.loc.cc.","rrsets":[
				{"name":"child.loc.cc.","type":"SOA","ttl":3600,"records":[{"content":"soa","disabled":false}]},
				{"name":"child.loc.cc.","type":"NS","ttl":3600,"records":[{"content":"ns1.example.","disabled":false}]},
				{"name":"child.loc.cc.","type":"A","ttl":300,"records":[{"content":"1.2.3.4","disabled":false}]},
				{"name":"www.child.loc.cc.","type":"CNAME","ttl":300,"records":[{"content":"child.loc.cc.","disabled":false}]}
			]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/servers/localhost/zones/loc.cc.":
			_, _ = w.Write([]byte(`{"name":"loc.cc.","rrsets":[
				{"name":"loc.cc.","type":"SOA","ttl":3600,"records":[{"content":"soa","disabled":false}]},
				{"name":"child.loc.cc.","type":"NS","ttl":3600,"records":[{"content":"ns1.example.","disabled":false}]},
				{"name":"child.loc.cc.","type":"DS","ttl":300,"records":[{"content":"1 13 2 AA","disabled":false}]},
				{"name":"old.child.loc.cc.","type":"A","ttl":300,"records":[{"content":"9.9.9.9","disabled":false}]},
				{"name":"other.loc.cc.","type":"A","ttl":300,"records":[{"content":"8.8.8.8","disabled":false}]}
			]}`))
		case r.Method == http.MethodPatch:
			var patch struct {
				RRsets []RRset `json:"rrsets"`
			}
			if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
				t.Fatalf("decode patch: %v", err)
			}
			patches[r.URL.Path] = append(patches[r.URL.Path], patch.RRsets...)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.SetDomainSuspended("loc.cc", "child.loc.cc", true); err != nil {
		t.Fatalf("SetDomainSuspended: %v", err)
	}

	parent := patches["/api/v1/servers/localhost/zones/loc.cc."]
	if len(parent) != 1 || parent[0].Name != "old.child.loc.cc." || !parent[0].Records[0].Disabled {
		t.Fatalf("parent: only the legacy record below the delegation should be disabled, got %+v", parent)
	}
	child := patches["/api/v1/servers/localhost/zones/child.loc.cc."]
	if len(child) != 2 {
		t.Fatalf("child: expected A and CNAME to be disabled, got %+v", child)
	}
	for _, rr := range child {
		if rr.Type == "SOA" || rr.Type == "NS" || !rr.Records[0].Disabled || rr.ChangeType != "REPLACE" {
			t.Fatalf("child: unexpected rrset %+v", rr)
		}
	}
}

func TestSetDomainSuspendedDisablesDelegationWhenNoChildZone(t *testing.T) {
	var parentPatch []RRset
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/servers/localhost/zones/ext.loc.cc.":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/servers/localhost/zones/loc.cc.":
			_, _ = w.Write([]byte(`{"name":"loc.cc.","rrsets":[
				{"name":"ext.loc.cc.","type":"NS","ttl":3600,"records":[{"content":"ns1.custom.","disabled":false}]},
				{"name":"ext.loc.cc.","type":"DS","ttl":300,"records":[{"content":"1 13 2 AA","disabled":false}]}
			]}`))
		case r.Method == http.MethodPatch:
			var patch struct {
				RRsets []RRset `json:"rrsets"`
			}
			_ = json.NewDecoder(r.Body).Decode(&patch)
			parentPatch = append(parentPatch, patch.RRsets...)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	if err := client.SetDomainSuspended("loc.cc", "ext.loc.cc", true); err != nil {
		t.Fatalf("SetDomainSuspended: %v", err)
	}
	if len(parentPatch) != 2 || !parentPatch[0].Records[0].Disabled || !parentPatch[1].Records[0].Disabled {
		t.Fatalf("custom-NS domain: NS and DS delegation must be disabled, got %+v", parentPatch)
	}
}

func TestRemoveDelegationsBatchesAndRectifiesOnce(t *testing.T) {
	var rrsetCount, rectifies int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			var patch struct {
				RRsets []RRset `json:"rrsets"`
			}
			_ = json.NewDecoder(r.Body).Decode(&patch)
			for _, rr := range patch.RRsets {
				if rr.ChangeType != "DELETE" || (rr.Type != "NS" && rr.Type != "DS") {
					t.Fatalf("unexpected rrset %+v", rr)
				}
			}
			rrsetCount += len(patch.RRsets)
		} else if r.Method == http.MethodPut {
			rectifies++
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	children := make([]string, 0, 70)
	for i := 0; i < 70; i++ {
		children = append(children, fmt.Sprintf("o%d.loc.cc", i))
	}
	if err := client.RemoveDelegations("loc.cc", children); err != nil {
		t.Fatalf("RemoveDelegations: %v", err)
	}
	if rrsetCount != 140 || rectifies != 1 {
		t.Fatalf("rrsets=%d rectifies=%d", rrsetCount, rectifies)
	}
}
