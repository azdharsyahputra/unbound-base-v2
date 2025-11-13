package user

import "time"

type Follow struct {
	ID          uint `gorm:"primaryKey"`
	FollowerID  uint `gorm:"not null;uniqueIndex:idx_follower_following"`
	FollowingID uint `gorm:"not null;uniqueIndex:idx_follower_following"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
