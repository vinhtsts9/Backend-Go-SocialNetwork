package service

import (
	"context"
	model "go-ecommerce-backend-api/m/v2/internal/models"
)

type (
	IPostProcessor interface {
		GetAllPosts(ctx context.Context) (codeResult int, posts []model.Post, err error)
		GetPostById(ctx context.Context, postId string) (codeResult int, post model.Post, err error)
		ProcessPostMessage(ctx context.Context)
	}
)

var localIPostProcessor IPostProcessor

func InitPostProcessor(i IPostProcessor) {
	localIPostProcessor = i
}

func PostProcessor() IPostProcessor {
	if localIPostProcessor == nil {
		panic("implement Post Processor not found")
	}
	return localIPostProcessor
}
