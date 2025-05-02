package model

type AppFeed struct {
	Type          string `json:"type"`
	AppIconURL    string `json:"app_icon_url"`
	AppBannerURL  string `json:"app_banner_url"`
	Version       string `json:"version"`
	Developer     string `json:"developer"`
	UpdatedOn     string `json:"updated_on"`
	DownloadCount string `json:"download_count"`
	ReleaseNotes  string `json:"release_notes"`
}
