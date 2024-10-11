//go:build wireinject

package wire

import (
	"go-ecommerce-backend-api/m/v2/internal/controller"
	"go-ecommerce-backend-api/m/v2/internal/repo"
	"go-ecommerce-backend-api/m/v2/internal/service"

	"github.com/google/wire"
)

func InitUserRouterHandler() (*controller.UserController, error) {
	wire.Build(
		repo.NewUserRepository,
		repo.NewUserAuthRepo,
		service.NewUserService,
		controller.NewUserController,
	)
	return new(controller.UserController), nil
}
