package router

type RouterGroup struct {
	BaseRouter
	AuthRouter
	OAuthRouter
	UserRouter
	ManageRouter
}

var RouterGroupApp = new(RouterGroup)
