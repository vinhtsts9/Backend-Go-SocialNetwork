package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
)

func InitStore() {

	store := database.NewStore(global.MdbcHaproxy)
	global.Store = store

}
