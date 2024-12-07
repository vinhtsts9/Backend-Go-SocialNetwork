package service

import (
	"context"
	model "go-ecommerce-backend-api/m/v2/internal/models"
)

type (
	IPost interface {
		CreatePost(ctx context.Context, input *model.CreatePostInput) (codeResult int, post model.Post, err error)
		UpdatePost(ctx context.Context, postId string, input *model.UpdatePostInput) (codeResult int, post model.Post, err error)
		DeletePost(ctx context.Context, postId string) (codeResult int, err error)
	}
)

var localPost IPost

func Post() IPost {
	if localPost == nil {
		panic("implement localPost notfound")
	}
	return localPost
}

func InitPost(i IPost) {
	localPost = i
}
