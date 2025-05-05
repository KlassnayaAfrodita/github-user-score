package refresh_cronjob

import (
	"context"
	github "github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/github-api"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/clients/repository"
	"github.com/KlassnayaAfrodita/github-user-score/collector/pkg/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"

	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/services"
	"github.com/robfig/cron/v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	db := database.NewDatabase(pool)
	repo := repository.NewCollectorRepository(db)
	githubClient := github.NewGitHubClient()

	collectorService := services.NewCollectorService(repo, githubClient)

	startCronJob(ctx, collectorService)

	select {}
}

func startCronJob(ctx context.Context, svc *services.CollectorService) {
	c := cron.New()

	_, err := c.AddFunc("@every 30m", func() {
		log.Println("[cron] starting outdated stats refresh")
		if err := svc.RefreshOutdatedStats(ctx, time.Hour); err != nil {
			log.Printf("[cron] refresh error: %v\n", err)
		} else {
			log.Println("[cron] refresh completed")
		}
	})

	if err != nil {
		log.Fatalf("failed to schedule cron job: %v", err)
	}

	c.Start()
}
