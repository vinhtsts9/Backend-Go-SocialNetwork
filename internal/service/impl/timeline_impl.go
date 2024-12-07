package impl

import (
	"fmt"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
)

type sTimeline struct {
	r *database.Queries
}

func NewTimelineImpl(r *database.Queries) *sTimeline {
	return &sTimeline{
		r: r,
	}
}

func (s *sTimeline) GetTimeLine(userId int) ([]model.Post, error) {
	key := fmt.Sprintf("timeline:%d", userId)

	timeline := getCache(key)
	return timeline, nil
}
