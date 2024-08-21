// scheduler/scheduler.go
package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type JobScheduler struct {
	redisClient *redis.Client
}

func NewJobScheduler(redisClient *redis.Client) *JobScheduler {
	return &JobScheduler{redisClient: redisClient}
}

func (js *JobScheduler) RunDailyJob(dailyJob func()) {
	dailyJobKey := "dailyJobRan"

	jobRan, err := js.redisClient.Exists(ctx, dailyJobKey).Result()
	if err != nil {
		log.Fatalf("Could not check Redis key: %v", err)
	}

	if jobRan == 0 {
		// If the job has not run today, execute it
		dailyJob()

		// Set the key in Redis with an expiration of 24 hours
		err := js.redisClient.Set(ctx, dailyJobKey, "true", 24*time.Hour).Err()
		if err != nil {
			log.Fatalf("Could not set Redis key: %v", err)
		}
		log.Println("Daily job completed.")
	} else {
		log.Println("Job already ran today. Skipping...")
	}
}

func (js *JobScheduler) Start() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			js.RunDailyJob(func() {
				log.Println("Running scheduled daily job...")
				// Your job logic here
			})
		}
	}
}
