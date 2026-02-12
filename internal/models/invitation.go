package models

import (
	"time"
)

type Invitation struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	InviterID   uint      `gorm:"not null" json:"inviter_id"`
	InviteeID   uint      `gorm:"not null" json:"invitee_id"`
	InviteCode  string    `gorm:"not null" json:"invite_code"`
	RewardGiven bool      `gorm:"default:false" json:"reward_given"`
	RewardType  string    `json:"reward_type"`
	RewardValue string    `json:"reward_value"`
	CreatedAt   time.Time `json:"created_at"`

	Inviter *User `gorm:"foreignKey:InviterID" json:"inviter,omitempty"`
	Invitee *User `gorm:"foreignKey:InviteeID" json:"invitee,omitempty"`
}

type InvitationResponse struct {
	ID          uint      `json:"id"`
	InviterID   uint      `json:"inviter_id"`
	InviteeID   uint      `json:"invitee_id"`
	InviteCode  string    `json:"invite_code"`
	RewardGiven bool      `json:"reward_given"`
	RewardType  string    `json:"reward_type"`
	RewardValue string    `json:"reward_value"`
	CreatedAt   time.Time `json:"created_at"`
	InviteeName string    `json:"invitee_name,omitempty"`
}

func (i *Invitation) ToResponse() *InvitationResponse {
	resp := &InvitationResponse{
		ID:          i.ID,
		InviterID:   i.InviterID,
		InviteeID:   i.InviteeID,
		InviteCode:  i.InviteCode,
		RewardGiven: i.RewardGiven,
		RewardType:  i.RewardType,
		RewardValue: i.RewardValue,
		CreatedAt:   i.CreatedAt,
	}

	if i.Invitee != nil {
		resp.InviteeName = i.Invitee.Username
	}

	return resp
}
