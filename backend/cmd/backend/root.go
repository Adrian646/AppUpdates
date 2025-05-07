package backend

import (
	"fmt"
	"github.com/Adrian646/AppUpdates/backend/internal/feeds/updater"
	"github.com/Adrian646/AppUpdates/backend/internal/handler"
	"github.com/Adrian646/AppUpdates/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func StartBackend() {
	fmt.Println("Starting backend...")

	err := InitDatabase(os.Getenv("DB_DSN"))
	if err != nil {
		return
	}

	updater.StartFeedUpdater(db)

	gin.SetMode(os.Getenv("GIN_MODE"))

	apiRoutePrefix := os.Getenv("API_ROUTE_PREFIX")

	r := gin.Default()

	r.Use(checkToken)

	r.GET(apiRoutePrefix+"feeds/updates", handler.GetFeedUpdates)
	r.GET(apiRoutePrefix+"feeds/:platform/:appID", handler.GetFeed)
	r.GET(apiRoutePrefix+"subscriptions/:subscriptionID", handler.GetSubscriptionByID)
	r.GET(apiRoutePrefix+"guilds/:guildID/feeds", handler.ListSubscriptions)
	r.POST(apiRoutePrefix+"guilds/:guildID/feeds", handler.CreateSubscription)
	r.DELETE(apiRoutePrefix+"guilds/:guildID/feeds/:platform/:appID", handler.DeleteSubscription)

	err = r.Run()

	if err != nil {
		log.Printf("Failed to start server: %v\n", err)
		return
	}
}

func checkToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != os.Getenv("API_SECRET") {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
}

func InitDatabase(dsn string) error {
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	if err := db.AutoMigrate(&model.AppFeed{}, &model.Subscription{}); err != nil {
		return fmt.Errorf("automigrate failed: %w", err)
	}

	handler.DB = db
	return nil
}
