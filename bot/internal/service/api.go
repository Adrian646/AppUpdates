package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type AppFeed struct {
	Platform      string    `json:"platform"`
	AppID         string    `json:"app_id"`
	AppName       string    `json:"app_name"`
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
	SubscriptionID uint    `json:"id"`
	GuildID        string  `json:"guild_id"`
	ChannelID      string  `json:"channel_id"`
	AppFeed        AppFeed `json:"app_feed"`
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

type Service struct {
	BaseURL string
	Client  *http.Client
	APIKey  string
}

func New(basePath string) *Service {
	host := os.Getenv("API_BASE_URL")
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	if strings.HasPrefix(basePath, "/") {
		basePath = strings.TrimPrefix(basePath, "/")
	}

	return &Service{
		BaseURL: host + basePath,
		Client:  &http.Client{Timeout: 1 * time.Minute},
		APIKey:  os.Getenv("API_SECRET"),
	}
}

func (s *Service) GetFeed(platform, appID string) (*AppFeed, error) {
	url := fmt.Sprintf("%sfeeds/%s/%s", s.BaseURL, platform, appID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.APIKey)

	resp, err := s.Client.Do(req)
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

func (s *Service) GetSubscriptionByID(subscriptionID string) (Subscription, error) {
	url := fmt.Sprintf("%ssubscriptions/%s", s.BaseURL, subscriptionID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Subscription{}, err
	}
	req.Header.Set("Authorization", s.APIKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		return Subscription{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Subscription{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var sub Subscription
	if err := json.NewDecoder(resp.Body).Decode(&sub); err != nil {
		return Subscription{}, err
	}
	return sub, nil
}

func (s *Service) ListSubscriptions(guildID string) ([]Subscription, error) {
	url := fmt.Sprintf("%sguilds/%s/feeds", s.BaseURL, guildID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.APIKey)

	resp, err := s.Client.Do(req)
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

func (s *Service) CreateSubscription(guildID, channelID, platform, appID string) (uint, error) {
	reqBody := map[string]string{
		"channel_id": channelID,
		"platform":   platform,
		"app_id":     appID,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}
	url := fmt.Sprintf("%sguilds/%s/feeds", s.BaseURL, guildID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.APIKey)

	resp, err := s.Client.Do(req)
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

func (s *Service) DeleteSubscription(guildID, platform, appID string) error {
	url := fmt.Sprintf("%sguilds/%s/feeds/%s/%s", s.BaseURL, guildID, platform, appID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", s.APIKey)

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

func (s *Service) GetFeedUpdates() ([]Subscription, error) {
	url := fmt.Sprintf("%sfeeds/updates", s.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.APIKey)

	resp, err := s.Client.Do(req)
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
