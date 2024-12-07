package user

type UserRouterGroup struct {
	UserRouter
	ProductRouter
	PostRouter
	RbacRouter
	ChatRouter
}
