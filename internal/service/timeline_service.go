package service

import model "go-ecommerce-backend-api/m/v2/internal/models"

type (
	TimelineInterface interface {
		GetTimeline(userId int) ([]model.Post, error)
	}
)

var localTimelineInterface TimelineInterface

func InitTimelineInterface(i TimelineInterface) {
	localTimelineInterface = i
}

func NewTimelineInterface() TimelineInterface {
	if localTimelineInterface == nil {
		panic("Failed to init timeline interface")
	}
	return localTimelineInterface
}
