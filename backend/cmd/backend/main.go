package main

import (
	"github.com/Adrian646/AppUpdates/backend/internal/feeds/android"
	"github.com/Adrian646/AppUpdates/backend/internal/feeds/ios"
)

func main() {
	android.GetCurrentAppData("io.bedrockhub.connector.bedrockhub_connector")
	_, err := ios.GetIOSAppData("6443529739")
	if err != nil {
		return
	}
}
