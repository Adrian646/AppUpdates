package model

import (
	"time"
)

type AppFeed struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Platform      string `gorm:"size:16;not null;index:idx_platform_appid,unique"`
	AppID         string `gorm:"size:64;not null;index:idx_platform_appid,unique"`
	Version       string `gorm:"size:32"`
	Developer     string `gorm:"size:128"`
	UpdatedOn     time.Time
	DownloadCount string
	AppIconURL    string    `gorm:"type:text"`
	AppBannerURL  string    `gorm:"type:text"`
	ReleaseNotes  string    `gorm:"type:text"`
	LastChecked   time.Time `gorm:"autoUpdateTime"`
	Notified      bool      `gorm:"not null;default:false;index"`
}

type Subscription struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	GuildID   string  `gorm:"size:32;not null;index"`
	ChannelID string  `gorm:"size:32;not null"`
	AppFeedID uint    `gorm:"not null;index"`
	AppFeed   AppFeed `gorm:"foreignKey:AppFeedID;constraint:OnDelete:CASCADE"`
}
