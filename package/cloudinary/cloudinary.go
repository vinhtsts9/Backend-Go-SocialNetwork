package cloudinary

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// Hàm khởi tạo Cloudinary
func InitCloudinary(cloudName, apiKey, apiSecret string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("Failed to init cloudinary")
	}
	return &CloudinaryService{
		cld: cld,
	}, nil
}

// Hàm upload ảnh từ URL
func (c *CloudinaryService) UploadImageFromURLToCloudinary(imageUrl string) (string, error) {
	// Kiểm tra Cloudinary đã được khởi tạo hay chưa
	if c.cld == nil {
		return "", errors.New("cloudinary not initialized")
	}

	// Upload từ URL
	resp, err := c.cld.Upload.Upload(context.Background(), imageUrl, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

// Hàm upload ảnh từ file
func (c *CloudinaryService) UploadImageToCloudinary(imagePath string) (string, error) {
	// Kiểm tra Cloudinary đã được khởi tạo hay chưa
	if c.cld == nil {
		return "", errors.New("cloudinary not initialized")
	}

	// Upload từ file
	resp, err := c.cld.Upload.Upload(context.Background(), imagePath, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
