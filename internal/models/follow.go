package models

import "time"

type Follow struct {
	FollowerID uint      `gorm:"primaryKey;column:follower_id"`
	FollowedID uint      `gorm:"primaryKey;column:followed_id"`
	FollowDate time.Time `gorm:"column:follow_date;default:CURRENT_TIMESTAMP"`
}