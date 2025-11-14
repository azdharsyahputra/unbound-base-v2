package model

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	User1ID   uint      `gorm:"not null" json:"user1_id"`
	User2ID   uint      `gorm:"not null" json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`

	Messages []Message `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
}

// Helper: return receiver user id
func (c *Chat) GetReceiver(senderID uint) uint {
	if c.User1ID == senderID {
		return c.User2ID
	}
	return c.User1ID
}
