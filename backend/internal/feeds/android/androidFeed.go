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

func GetCurrentAppData(appID string) model.AppFeed {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&hl=en&gl=US", appID)

	var labels []string
	var values []string
	var iconURL string
	var screenshotURL string
	var releaseNotes string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`img.T75of`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return clickSeeMoreIfPresent(ctx)
		}),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Evaluate(`Array.from(document.querySelectorAll("div.G1zzid .q078ud")).map(e => e.textContent)`, &labels),
		chromedp.Evaluate(`Array.from(document.querySelectorAll("div.G1zzid .reAt0")).map(e => e.textContent)`, &values),
		chromedp.AttributeValue(`img.T75of`, "src", &iconURL, nil, chromedp.ByQuery),
		chromedp.AttributeValue(`img.T75of.B5GQxf`, "src", &screenshotURL, nil, chromedp.ByQuery),
		chromedp.Text(`div[itemprop="description"]`, &releaseNotes, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatalf("chromedp failed: %v", err)
	}

	appFeed := model.AppFeed{
		Type:         "android",
		AppIconURL:   iconURL,
		AppBannerURL: screenshotURL,
		ReleaseNotes: releaseNotes,
	}

	for i := 0; i < len(labels) && i < len(values); i++ {
		label := strings.TrimSpace(labels[i])
		value := strings.TrimSpace(values[i])

		switch label {
		case "Version":
			appFeed.Version = value
		case "Updated on":
			appFeed.UpdatedOn = value
		case "Offered by":
			appFeed.Developer = value
		case "Downloads":
			if idx := strings.Index(value, "+"); idx != -1 {
				appFeed.DownloadCount = strings.TrimSpace(value[:idx])
			} else {
				appFeed.DownloadCount = value
			}
		}

	}

	fmt.Printf("%+v\n", appFeed)
	return appFeed
}

func clickSeeMoreIfPresent(ctx context.Context) error {
	selectors := []string{
		`button[aria-label="See more information on About this app"]`,
		`button[aria-label="See more information on About this game"]`,
	}

	for _, sel := range selectors {
		var visible bool
		err := chromedp.Run(ctx,
			chromedp.EvaluateAsDevTools(fmt.Sprintf(`!!document.querySelector('%s')`, sel), &visible),
		)
		if err == nil && visible {
			return chromedp.Run(ctx,
				chromedp.ScrollIntoView(sel, chromedp.ByQuery),
				chromedp.Click(sel, chromedp.ByQuery),
			)
		}
	}

	log.Println("No button found to read the app information's.")
	return nil
}
