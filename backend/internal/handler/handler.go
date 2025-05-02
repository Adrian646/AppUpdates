package handler

import (
	"errors"
	"github.com/Adrian646/AppUpdates/backend/internal/feeds/android"
	"github.com/Adrian646/AppUpdates/backend/internal/feeds/ios"
	"github.com/Adrian646/AppUpdates/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var DB *gorm.DB

const feedTTL = 15 * time.Minute

func GetFeed(c *gin.Context) {
	platform := c.Param("platform")
	appID := c.Param("appID")

	var feed model.AppFeed
	err := DB.
		Where("platform = ? AND app_id = ?", platform, appID).
		First(&feed).Error

	miss := errors.Is(err, gorm.ErrRecordNotFound)
	stale := !miss && time.Since(feed.LastChecked) > feedTTL

	if miss || stale {
		var fresh model.AppFeed
		if platform == "ios" {
			fresh, err = ios.GetCurrentAppData(appID)
		} else {
			fresh, err = android.GetCurrentAppData(appID)
		}
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		feed.Platform = platform
		feed.AppID = appID
		feed.Version = fresh.Version
		feed.Developer = fresh.Developer
		feed.UpdatedOn = fresh.UpdatedOn
		feed.DownloadCount = fresh.DownloadCount
		feed.AppIconURL = fresh.AppIconURL
		feed.AppBannerURL = fresh.AppBannerURL
		feed.ReleaseNotes = fresh.ReleaseNotes
		feed.LastChecked = time.Now()

		if miss {
			if err := DB.Create(&feed).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			if err := DB.Save(&feed).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"platform":       feed.Platform,
		"app_id":         feed.AppID,
		"version":        feed.Version,
		"developer":      feed.Developer,
		"updated_on":     feed.UpdatedOn,
		"download_count": feed.DownloadCount,
		"app_icon_url":   feed.AppIconURL,
		"app_banner_url": feed.AppBannerURL,
		"release_notes":  feed.ReleaseNotes,
	})
}

func ListSubscriptions(c *gin.Context) {
	guildID := c.Param("guildID")
	var subs []model.Subscription

	if err := DB.
		Preload("AppFeed").
		Where("guild_id = ?", guildID).
		Find(&subs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]gin.H, len(subs))
	for i, s := range subs {
		out[i] = gin.H{
			"subscription_id": s.ID,
			"channel_id":      s.ChannelID,
			"platform":        s.AppFeed.Platform,
			"app_id":          s.AppFeed.AppID,
			"feed": gin.H{
				"version":        s.AppFeed.Version,
				"developer":      s.AppFeed.Developer,
				"updated_on":     s.AppFeed.UpdatedOn,
				"download_count": s.AppFeed.DownloadCount,
				"app_icon_url":   s.AppFeed.AppIconURL,
				"app_banner_url": s.AppFeed.AppBannerURL,
				"release_notes":  s.AppFeed.ReleaseNotes,
			},
		}
	}

	c.JSON(http.StatusOK, out)
}

func CreateSubscription(c *gin.Context) {
	guildID := c.Param("guildID")
	var req struct {
		ChannelID string `json:"channel_id" binding:"required"`
		Platform  string `json:"platform"  binding:"required,oneof=ios android"`
		AppID     string `json:"app_id"    binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed := model.AppFeed{Platform: req.Platform, AppID: req.AppID}

	if err := DB.Where("platform = ? AND app_id = ?", req.Platform, req.AppID).
		FirstOrCreate(&feed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sub := model.Subscription{
		GuildID:   guildID,
		ChannelID: req.ChannelID,
		AppFeedID: feed.ID,
	}
	if err := DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription_id": sub.ID, "message": "Subscription created"})
}

func DeleteSubscription(c *gin.Context) {
	guildID := c.Param("guildID")
	platform := c.Param("platform")
	appID := c.Param("appID")

	var feed model.AppFeed
	if err := DB.
		Where("platform = ? AND app_id = ?", platform, appID).
		First(&feed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No feed found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	result := DB.
		Where("guild_id = ? AND app_feed_id = ?", guildID, feed.ID).
		Delete(&model.Subscription{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No subscription found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}
