package updater

import (
	"log"
	"time"

	"github.com/Adrian646/AppUpdates/backend/internal/feeds/android"
	"github.com/Adrian646/AppUpdates/backend/internal/feeds/ios"
	"github.com/Adrian646/AppUpdates/backend/internal/model"
	"gorm.io/gorm"
)

const feedTTL = 3 * time.Minute

func StartFeedUpdater(db *gorm.DB) {
	ticker := time.NewTicker(feedTTL)
	go func() {
		for range ticker.C {
			checkAllFeeds(db)
		}
	}()
}

func checkAllFeeds(db *gorm.DB) {
	var feeds []model.AppFeed
	if err := db.Find(&feeds).Error; err != nil {
		log.Printf("Couldnt retrieve any feeds from database: %v\n", err)
		return
	}

	for _, f := range feeds {
		var fresh model.AppFeed
		var err error
		switch f.Platform {
		case "ios":
			fresh, err = ios.GetCurrentAppData(f.AppID)
		case "android":
			fresh, err = android.GetCurrentAppData(f.AppID)
		default:
			continue
		}
		if err != nil {
			log.Printf("Error while retrieving feed: %s/%s: %v\n", f.Platform, f.AppID, err)
			continue
		}

		if fresh.Version != f.Version {
			f.Version = fresh.Version
			f.UpdatedOn = fresh.UpdatedOn
			f.ReleaseNotes = fresh.ReleaseNotes
			f.DownloadCount = fresh.DownloadCount
			f.AppIconURL = fresh.AppIconURL
			f.AppBannerURL = fresh.AppBannerURL
			f.Developer = fresh.Developer
			f.Notified = false
			if err := db.Save(&f).Error; err != nil {
				log.Printf("Couldnt save the newest feed: %v\n", err)
			} else {
				log.Printf("Found a new version: %s/%s → %s\n", f.Platform, f.AppName, f.Version)
			}
		} else {
			db.Model(&f).Update("last_checked", time.Now())
		}
	}
}
