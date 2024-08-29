package scheduler

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type JobScheduler struct {
	redisClient    *redis.Client
	github_id      string
	githubUsername string
}

func NewJobScheduler(redisClient *redis.Client, githubID, githubUsername string) *JobScheduler {
	return &JobScheduler{
		redisClient:    redisClient,
		github_id:      githubID,
		githubUsername: githubUsername,
	}
}

func (js *JobScheduler) RunDailyJob(dailyJob func()) {
	dailyJobKey := fmt.Sprintf("dailyJobRan:%s", js.github_id)

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
				js.PerformGithubCommit()
			})
		}
	}
}

func (js *JobScheduler) PerformGithubCommit() {
	log.Printf("Performing daily commit for github User: %s\n", js.githubUsername)

	commitMsg := fmt.Sprintf("Daily commit by %s", js.githubUsername)

	addCmd := exec.Command("git", "add", ".")
	err := addCmd.Run()
	if err !=  nil {
		log.Printf("Error adding changes: %v\n", err)
		return
	}

	commitCmd := exec.Command("git","commit", "-m", commitMsg)
	err = commitCmd.Run()
	if err != nil {
		log.Printf("Error committing changes: %v\n", err)
		return
	}

	pushCmd := exec.Command("git", "push")
	err = pushCmd.Run()
	if err != nil {
		log.Printf("Error pushing changes: %v\n", err)
		return
	}

	log.Println("Daily commit completed successfully ")
}

func TestJob() {
	log.Println("Running test job...")

	time.Sleep((1 * time.Second))
	log.Println("Test job completed.")
}
