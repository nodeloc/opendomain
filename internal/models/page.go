package models

import "time"

type Page struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string    `gorm:"size:200;not null" json:"title"`
	Slug         string    `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Content      string    `gorm:"type:text;not null" json:"content"`
	Category     string    `gorm:"size:50;index;not null" json:"category"` // company, resources
	IsPublished  bool      `gorm:"default:true" json:"is_published"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Page) TableName() string {
	return "pages"
}
