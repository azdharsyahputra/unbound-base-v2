package notification

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	SenderID uint   `json:"sender_id"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Read     bool   `json:"read" gorm:"default:false"`
}

func (Notification) TableName() string {
	return "notifications"
}
