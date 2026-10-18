package models

import "time"

const (
	SponsorTierPartner = "partner"
	SponsorTierSponsor = "sponsor"
	SponsorTierBacker  = "backer"

	SettingSponsorUSDTTRC20Address = "sponsor_usdt_trc20_address"
	SettingSponsorGithubURL        = "sponsor_github_url"
)

// SponsorTiers is listed in public display order.
var SponsorTiers = []string{SponsorTierPartner, SponsorTierSponsor, SponsorTierBacker}

// Sponsor 赞助商
//
// IsActive and DisplayOrder intentionally carry no GORM default tag: with a
// default tag GORM omits zero values (false / 0) on insert, so they would
// silently fall back to the database default.
type Sponsor struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	WebsiteURL    string    `gorm:"column:website_url;size:500;not null" json:"website_url"`
	LogoURL       string    `gorm:"column:logo_url;size:500;not null" json:"logo_url"`
	LogoDarkURL   string    `gorm:"column:logo_dark_url;size:500;not null" json:"logo_dark_url"`
	DescriptionZh string    `gorm:"column:description_zh;size:300;not null" json:"description_zh"`
	DescriptionEn string    `gorm:"column:description_en;size:300;not null" json:"description_en"`
	Tier          string    `gorm:"size:20;not null" json:"tier"`
	DisplayOrder  int       `gorm:"not null" json:"display_order"`
	IsActive      bool      `gorm:"not null" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Sponsor) TableName() string {
	return "sponsors"
}

// SponsorRequest 管理员创建/更新赞助商请求（PUT 为整体替换）
type SponsorRequest struct {
	Name          string `json:"name"`
	WebsiteURL    string `json:"website_url"`
	LogoURL       string `json:"logo_url"`
	LogoDarkURL   string `json:"logo_dark_url"`
	DescriptionZh string `json:"description_zh"`
	DescriptionEn string `json:"description_en"`
	Tier          string `json:"tier"`
	DisplayOrder  int    `json:"display_order"`
	IsActive      *bool  `json:"is_active"`
}

// PublicSponsor 公开赞助商信息
type PublicSponsor struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	WebsiteURL    string `json:"website_url"`
	LogoURL       string `json:"logo_url"`
	LogoDarkURL   string `json:"logo_dark_url"`
	DescriptionZh string `json:"description_zh"`
	DescriptionEn string `json:"description_en"`
	Tier          string `json:"tier"`
}

func (s *Sponsor) ToPublic() PublicSponsor {
	return PublicSponsor{
		ID:            s.ID,
		Name:          s.Name,
		WebsiteURL:    s.WebsiteURL,
		LogoURL:       s.LogoURL,
		LogoDarkURL:   s.LogoDarkURL,
		DescriptionZh: s.DescriptionZh,
		DescriptionEn: s.DescriptionEn,
		Tier:          s.Tier,
	}
}

// SponsorSettings 赞助付款方式设置
type SponsorSettings struct {
	USDTTRC20Address  string `json:"usdt_trc20_address"`
	GithubSponsorsURL string `json:"github_sponsors_url"`
}
