package android

import (
	"fmt"
	"github.com/Adrian646/AppUpdates/backend/internal/model"
	"github.com/n0madic/google-play-scraper/pkg/app"
)

func GetCurrentAppData(appID string) (model.AppFeed, error) {
	a := app.New(appID, app.Options{
		Country:  "us",
		Language: "us",
	})

	err := a.LoadDetails()
	if err != nil {
		return model.AppFeed{}, err
	}

	feed := model.AppFeed{
		Platform:      "android",
		AppID:         a.ID,
		AppName:       a.Title,
		AppIconURL:    a.Icon,
		AppBannerURL:  a.Screenshots[0],
		ReleaseNotes:  a.RecentChanges,
		Version:       a.Version,
		Developer:     a.Developer,
		DownloadCount: a.Installs,
		UpdatedOn:     a.Updated,
	}

	fmt.Println(feed)

	return feed, nil
}
