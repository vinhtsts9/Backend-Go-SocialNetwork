package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
)

type sPostProcessor struct {
	r *database.Queries
}

func NewPostProcessorImpl(r *database.Queries) *sPostProcessor {
	return &sPostProcessor{
		r: r,
	}
}

// Cập nhật cache timeline
func UpdateTimelineCache(post model.Post, onlineFollowers []model.UserBase) error {
	const batchsize = 100
	const maxWorkers = 5

	jobs := make(chan []model.UserBase, maxWorkers)
	errChan := make(chan error, maxWorkers)
	worker := func() {
		for batch := range jobs {
			pipeline := global.Rdb.Pipeline()
			for _, follower := range batch {
				key := fmt.Sprintf("timeline%d", follower.UserID)

				timeline := getCache(key)
				timeline = append(timeline, post)

				data, err := json.Marshal(timeline)
				if err != nil {
					errChan <- fmt.Errorf("failed to marshal timeline %w", err)
					return
				}
				pipeline.Set(context.Background(), key, data, 0)
			}

			_, err := pipeline.Exec(context.Background())
			if err != nil {
				errChan <- fmt.Errorf("pipeline execution failed: %w", err)
				return
			}
		}
		errChan <- nil // bao thanh cong khi worker hoan thanh
	}
	for i := 0; i < maxWorkers; i++ {
		go worker()
	}
	go func() {
		for i := 0; i < len(onlineFollowers); i += batchsize {
			end := i + batchsize
			if end > len(onlineFollowers) {
				end = len(onlineFollowers)
			}
			batch := onlineFollowers[i:end]
			jobs <- batch
		}
		close(jobs)
	}()

	for i := 0; i < maxWorkers; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}

// Lấy cache từ Redis
func getCache(key string) []model.Post {
	var timeline []model.Post
	data, err := global.Rdb.Get(context.Background(), key).Result()
	if err == nil {
		json.Unmarshal([]byte(data), &timeline)
	}
	return timeline
}

// Ghi cache vào Redis
func setCache(key string, timeline *[]model.Post) error {
	data, err := json.Marshal(&timeline)
	if err != nil {
		return fmt.Errorf("failed to marshal timeline: %w", err)
	}
	return global.Rdb.SetEx(context.Background(), key, data, 10).Err()
}
