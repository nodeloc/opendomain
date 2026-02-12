package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Type        string         `gorm:"not null" json:"type"` // general, maintenance, update, important
	Priority    int            `gorm:"default:0" json:"priority"`
	IsPublished bool           `gorm:"default:false" json:"is_published"`
	PublishedAt *time.Time     `json:"published_at"`
	AuthorID    *uint          `json:"author_id"`
	Views       int            `gorm:"default:0" json:"views"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Author *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

type AnnouncementCreateRequest struct {
	Title    string `json:"title" binding:"required,min=3,max=255"`
	Content  string `json:"content" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=general maintenance update important"`
	Priority int    `json:"priority"`
}

type AnnouncementUpdateRequest struct {
	Title       *string `json:"title"`
	Content     *string `json:"content"`
	Type        *string `json:"type"`
	Priority    *int    `json:"priority"`
	IsPublished *bool   `json:"is_published"`
}

type AnnouncementResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Type        string     `json:"type"`
	Priority    int        `json:"priority"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
	Views       int        `json:"views"`
	CreatedAt   time.Time  `json:"created_at"`
	AuthorName  string     `json:"author_name,omitempty"`
}

func (a *Announcement) ToResponse() *AnnouncementResponse {
	resp := &AnnouncementResponse{
		ID:          a.ID,
		Title:       a.Title,
		Content:     a.Content,
		Type:        a.Type,
		Priority:    a.Priority,
		IsPublished: a.IsPublished,
		PublishedAt: a.PublishedAt,
		Views:       a.Views,
		CreatedAt:   a.CreatedAt,
	}

	if a.Author != nil {
		resp.AuthorName = a.Author.Username
	}

	return resp
}
