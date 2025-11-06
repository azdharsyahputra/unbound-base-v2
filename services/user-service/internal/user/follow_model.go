package user

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	FollowerID  uint `gorm:"not null;uniqueIndex:idx_follower_following"`
	FollowingID uint `gorm:"not null;uniqueIndex:idx_follower_following"`
}
