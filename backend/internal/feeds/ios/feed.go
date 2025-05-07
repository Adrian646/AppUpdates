package ios

import (
	"encoding/json"
	"fmt"
	"github.com/Adrian646/AppUpdates/backend/internal/model"
	"log"
	"net/http"
	"time"
)

type iOSLookupResponse struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		TrackName          string   `json:"trackName"`
		Version            string   `json:"version"`
		ArtistName         string   `json:"artistName"`
		ArtworkUrl512      string   `json:"artworkUrl512"`
		ScreenshotUrls     []string `json:"screenshotUrls"`
		CurrentVersionDate string   `json:"currentVersionReleaseDate"`
		ReleaseNotes       string   `json:"releaseNotes"`
	} `json:"results"`
}

func GetCurrentAppData(appID string) (model.AppFeed, error) {
	var feed model.AppFeed

	url := fmt.Sprintf("https://itunes.apple.com/lookup?id=%s&country=US", appID)
	client := &http.Client{Timeout: 20 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return feed, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "AppUpdatesBot/1.0")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return feed, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return feed, fmt.Errorf("non-200 status: %d", resp.StatusCode)
	}

	var apiResp iOSLookupResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return feed, fmt.Errorf("error decoding response: %w", err)
	}
	if apiResp.ResultCount == 0 || len(apiResp.Results) == 0 {
		return feed, fmt.Errorf("no app found for ID %s", appID)
	}
	app := apiResp.Results[0]

	updatedOn, err := time.Parse(time.RFC3339, app.CurrentVersionDate)
	if err != nil {
		return feed, fmt.Errorf("error parsing date %q: %w", app.CurrentVersionDate, err)
	}

	feed.Platform = "ios"
	feed.AppID = appID
	feed.AppName = app.TrackName
	feed.Version = app.Version
	feed.Developer = app.ArtistName
	feed.AppIconURL = app.ArtworkUrl512
	if len(app.ScreenshotUrls) > 0 {
		feed.AppBannerURL = app.ScreenshotUrls[0]
	}
	feed.ReleaseNotes = app.ReleaseNotes
	feed.UpdatedOn = updatedOn

	log.Printf("[iOS] AppID %s → Version %s (%s)", appID, app.Version, feed.UpdatedOn.Format("2006-01-02"))

	return feed, nil
}
