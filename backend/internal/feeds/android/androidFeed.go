package android

import (
	"context"
	"fmt"
	"github.com/Adrian646/AppUpdates/backend/internal/model"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

func GetCurrentAppData(appID string) (model.AppFeed, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&hl=en&gl=US", appID)

	var (
		labels        []string
		values        []string
		iconURL       string
		screenshotURL string
		releaseNotes  string
	)

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`img.T75of`, chromedp.ByQuery),
		chromedp.ActionFunc(clickSeeMoreIfPresent),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Evaluate(`Array.from(document.querySelectorAll("div.G1zzid .q078ud")).map(e => e.textContent)`, &labels),
		chromedp.Evaluate(`Array.from(document.querySelectorAll("div.G1zzid .reAt0")).map(e => e.textContent)`, &values),
		chromedp.AttributeValue(`img.T75of`, "src", &iconURL, nil),
		chromedp.AttributeValue(`img.T75of.B5GQxf`, "src", &screenshotURL, nil),
		chromedp.Text(`div[itemprop="description"]`, &releaseNotes, chromedp.ByQuery),
	)
	if err != nil {
		return model.AppFeed{}, fmt.Errorf("chromedp failed: %w", err)
	}

	feed := model.AppFeed{
		Platform:     "android",
		AppID:        appID,
		AppIconURL:   iconURL,
		AppBannerURL: screenshotURL,
		ReleaseNotes: releaseNotes,
	}

	print("a")

	for i := 0; i < len(labels) && i < len(values); i++ {
		label := strings.TrimSpace(labels[i])
		value := strings.TrimSpace(values[i])

		switch label {
		case "Version":
			feed.Version = value

		case "Updated on":
			fmt.Println(value)
			if t, err := time.Parse("Jan 2, 2006", value); err == nil {
				feed.UpdatedOn = t
			}

		case "Offered by":
			feed.Developer = value

		case "Downloads":
			if idx := strings.Index(value, "+"); idx != -1 {
				feed.DownloadCount = strings.TrimSpace(value[:idx])
			} else {
				feed.DownloadCount = value
			}
		}
	}

	return feed, nil
}

func clickSeeMoreIfPresent(ctx context.Context) error {
	selectors := []string{
		`button[aria-label="See more information on About this app"]`,
		`button[aria-label="See more information on About this game"]`,
	}
	for _, sel := range selectors {
		var visible bool
		if err := chromedp.Run(ctx,
			chromedp.EvaluateAsDevTools(fmt.Sprintf(`!!document.querySelector('%s')`, sel), &visible),
		); err == nil && visible {
			if err := chromedp.Run(ctx,
				chromedp.ScrollIntoView(sel, chromedp.ByQuery),
				chromedp.Click(sel, chromedp.ByQuery),
			); err != nil {
				log.Println("clickSeeMore error:", err)
			}
			break
		}
	}
	return nil
}
