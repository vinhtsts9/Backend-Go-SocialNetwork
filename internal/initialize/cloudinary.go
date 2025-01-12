package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/package/cloudinary"
	"os"
)

func NewCloudinary() {
	cloud_name := os.Getenv("CLOUD_NAME")
	api_key := os.Getenv("API_KEY")
	api_secret := os.Getenv("API_SECRET")
	global.Logger.Sugar().Info(cloud_name, api_key, api_secret)
	Cloudinary, err := cloudinary.InitCloudinary(cloud_name, api_key, api_secret)
	if err != nil {

	}
	global.Cloudinary = Cloudinary
}
