package model

type AppFeed struct {
	AppIconURL    string `json:"app_icon_url"`
	AppBannerURL  string `json:"app_banner_url"`
	Version       string `json:"version"`
	Developer     string `json:"developer"`
	UpdatedOn     string `json:"updated_on"`
	DownloadCount string `json:"download_count"`
}
