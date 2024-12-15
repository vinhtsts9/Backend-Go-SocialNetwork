package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/package/cloudinary"
)

func NewCloudinary() {
	cloud_name := "vinhts"
	api_key := "194494324816216"
	api_secret := "qpmox9UBhnfx6I_gi-rO_eR2eRA"
	Cloudinary, err := cloudinary.InitCloudinary(cloud_name, api_key, api_secret)
	if err != nil {

	}
	global.Cloudinary = Cloudinary
}
