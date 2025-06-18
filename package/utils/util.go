package utils

import (
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetUserKey(hashKey string) string {
	return fmt.Sprintf("u:%s:otp", hashKey)
}
func GenerateCliTokenUUID(userId int) string {
	newUUID := uuid.New()
	uuidString := strings.ReplaceAll((newUUID).String(), "", "")
	return strconv.Itoa(userId) + "clitoken" + uuidString
}
func MapGetPostByIdRowToPost(row database.GetPostByIdRow) (model.Post, error) {
	// Giải mã ImagePaths trực tiếp
	var imagePaths []string
	if err := json.Unmarshal(row.ImagePaths, &imagePaths); err != nil {
		return model.Post{}, fmt.Errorf("failed to unmarshal image paths: %w", err)
	}

	// Xử lý CreatedAt và UpdatedAt
	createdAt := ""
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time.Format(time.RFC3339)
	}
	updatedAt := ""
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time.Format(time.RFC3339)
	}

	return model.Post{
		ID:           uint32(row.ID),
		UserId:       0, // Cần thêm nếu có trường UserId
		Title:        row.Title,
		ImagePaths:   imagePaths,
		UserNickname: row.UserNickname,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		IsPublished:  false, // Cần thêm nếu có trường IsPublished
		Metadata:     "",    // Xử lý nếu Metadata là JSON
	}, nil
}

func MapGetAllpostRowToPost(row database.GetAllpostRow) (model.Post, error) {
	// Giải mã ImagePaths trực tiếp
	var imagePaths []string
	if err := json.Unmarshal(row.ImagePaths, &imagePaths); err != nil {
		return model.Post{}, fmt.Errorf("failed to unmarshal image paths: %w", err)
	}

	// Xử lý CreatedAt và UpdatedAt
	createdAt := ""
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time.Format(time.RFC3339)
	}
	updatedAt := ""
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time.Format(time.RFC3339)
	}

	return model.Post{
		ID:           uint32(row.ID),
		Title:        row.Title,
		ImagePaths:   imagePaths,
		UserNickname: row.UserNickname,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
