package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/package/cloudinary"
)

func NewCloudinary() {
	cloud_name := ""
	api_key := ""
	api_secret := ""
	Cloudinary, err := cloudinary.InitCloudinary(cloud_name, api_key, api_secret)
	if err != nil {

	}
	global.Cloudinary = Cloudinary
}
