package initialize

import (
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()
	viper.AddConfigPath("./configs")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config %w", err))
	}
	viper.AddConfigPath(".")   // Đặt đường dẫn tìm file .env
	viper.SetConfigName("app") // Tên file .env (app.env)
	viper.SetConfigType("env") // Định dạng .env

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config from app.env:", err)
	}

	viper.AutomaticEnv()

	fmt.Println("server port", viper.GetInt("server.port"))
	fmt.Println("security port", viper.GetString("security.jwt.key"))

	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unable to decode configuration %v", err)
	}
}
