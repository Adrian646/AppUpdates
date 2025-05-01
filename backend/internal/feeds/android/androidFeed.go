package android

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func GetCurrentAppData() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := buildURL("com.tencent.ig")

	var labels []string
	var values []string
	var iconURL string
	var screenshotURL string

	var buttonFound bool

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`img.T75of`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := chromedp.Click(`button[aria-label="See more information on About this app"]`, chromedp.ByQuery).Do(ctx)
			if err == nil {
				buttonFound = true
				return nil
			}

			err = chromedp.Click(`button[aria-label="See more information on About this game"]`, chromedp.ByQuery).Do(ctx)
			if err == nil {
				buttonFound = true
			}
			return nil
		}),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Evaluate(`Array.from(document.querySelectorAll("div.G1zzid .q078ud")).map(e => e.textContent)`, &labels),
		chromedp.Evaluate(`Array.from(document.querySelectorAll("div.G1zzid .reAt0")).map(e => e.textContent)`, &values),
		chromedp.AttributeValue(`img.T75of`, "src", &iconURL, nil, chromedp.ByQuery),
		chromedp.AttributeValue(`img.T75of.B5GQxf`, "src", &screenshotURL, nil, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatalf("chromedp failed: %v", err)
	}

	if !buttonFound {
		log.Println("Kein 'See more information'-Button gefunden.")
	}

	for i := 0; i < len(labels) && i < len(values); i++ {
		fmt.Printf("%s: %s\n", labels[i], values[i])
	}

	fmt.Printf("Icon URL: %s\n", iconURL)
	fmt.Printf("Screenshot URL: %s\n", screenshotURL)
}

func buildURL(appId string) string {
	return "https://play.google.com/store/apps/details?id=" + appId + "&hl=en&gl=US"
}
