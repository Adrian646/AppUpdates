package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AppFeed struct {
	Platform      string    `json:"platform"`
	AppID         string    `json:"app_id"`
	Version       string    `json:"version"`
	Developer     string    `json:"developer"`
	UpdatedOn     time.Time `json:"updated_on"`
	DownloadCount string    `json:"download_count"`
	AppIconURL    string    `json:"app_icon_url"`
	AppBannerURL  string    `json:"app_banner_url"`
	ReleaseNotes  string    `json:"release_notes"`
	NewVersion    bool      `json:"new_version,omitempty"`
}

type Subscription struct {
	SubscriptionID uint    `json:"subscription_id"`
	ChannelID      string  `json:"channel_id"`
	Platform       string  `json:"platform"`
	AppID          string  `json:"app_id"`
	Feed           AppFeed `json:"feed"`
}

type GuildUpdate struct {
	SubscriptionID uint      `json:"subscription_id"`
	Platform       string    `json:"platform"`
	AppID          string    `json:"app_id"`
	Version        string    `json:"version"`
	Developer      string    `json:"developer"`
	UpdatedOn      time.Time `json:"updated_on"`
	DownloadCount  string    `json:"download_count"`
	AppIconURL     string    `json:"app_icon_url"`
	AppBannerURL   string    `json:"app_banner_url"`
	ReleaseNotes   string    `json:"release_notes"`
}

type APIService struct {
	BaseURL string
	Client  *http.Client
}

func New(baseURL string) *APIService {
	return &APIService{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *APIService) GetFeed(platform, appID string) (*AppFeed, error) {
	url := fmt.Sprintf("%s/api/v1/feeds/%s/%s", s.BaseURL, platform, appID)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var feed AppFeed
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}
	return &feed, nil
}

func (s *APIService) ListSubscriptions(guildID string) ([]Subscription, error) {
	url := fmt.Sprintf("%s/api/v1/guilds/%s/feeds", s.BaseURL, guildID)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var subs []Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *APIService) CreateSubscription(guildID, channelID, platform, appID string) (uint, error) {
	reqBody := map[string]string{
		"channel_id": channelID,
		"platform":   platform,
		"app_id":     appID,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}
	url := fmt.Sprintf("%s/api/v1/guilds/%s/feeds", s.BaseURL, guildID)
	resp, err := s.Client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var result struct {
		SubscriptionID uint   `json:"subscription_id"`
		Message        string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.SubscriptionID, nil
}

func (s *APIService) DeleteSubscription(guildID, platform, appID string) error {
	url := fmt.Sprintf("%s/api/v1/guilds/%s/feeds/%s/%s", s.BaseURL, guildID, platform, appID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

func (s *APIService) GetGuildUpdates(guildID string) ([]GuildUpdate, error) {
	url := fmt.Sprintf("%s/api/v1/guilds/%s/updates", s.BaseURL, guildID)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var msg struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &msg); err == nil && msg.Message != "" {
		return []GuildUpdate{}, nil
	}
	var updates []GuildUpdate
	if err := json.Unmarshal(body, &updates); err != nil {
		return nil, err
	}
	return updates, nil
}
