package model

import "time"

type Message struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ChatID   uint   `gorm:"index;not null" json:"chat_id"`
	SenderID uint   `gorm:"not null" json:"sender_id"`
	Content  string `gorm:"type:text;not null" json:"content"`

	Status string     `gorm:"type:varchar(20);default:'sent'" json:"status"`
	IsRead bool       `gorm:"default:false" json:"is_read"`
	ReadAt *time.Time `json:"read_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
