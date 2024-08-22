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
		dailyJob()

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

func TestJob(){
	log.Println("Running test job...")
	//simulate work
	time.Sleep((1 * time.Second))
	log.Println("Test job completed.")
}

func (js *JobScheduler) TriggerTestJob() {
	js.RunDailyJob(TestJob)
}