package routers

import (
	"go-ecommerce-backend-api/m/v2/internal/routers/manage"
	"go-ecommerce-backend-api/m/v2/internal/routers/user"
)

type RouterGroup struct {
	User   user.UserRouterGroup
	Manage manage.ManageRouterGroup
}

var RouterGroupApp = new(RouterGroup)
