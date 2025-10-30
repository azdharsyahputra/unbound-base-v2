package user

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	FollowerID  uint `gorm:"not null"`
	FollowingID uint `gorm:"not null"`
}
