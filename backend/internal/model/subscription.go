package model

import "time"

type Subscription struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	GuildID   string  `gorm:"size:32;not null;index"`
	ChannelID string  `gorm:"size:32;not null"`
	AppFeedID uint    `gorm:"not null;index"`
	AppFeed   AppFeed `gorm:"foreignKey:AppFeedID;constraint:OnDelete:CASCADE"`
}
