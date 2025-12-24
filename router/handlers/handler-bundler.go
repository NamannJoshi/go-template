package handlers

import "template/services"

type RouteHandlers struct {
	UserRoutes IUserHandler
	AuthRoutes IAuthHandler
}

func RouteBundler(svc services.Services) *RouteHandlers {
	return &RouteHandlers{
		UserRoutes: UserHandlerInit(svc.UserServices, svc.CredServices),
		AuthRoutes: AuthHandlerInit(svc.AuthServices, svc.CredServices),
	}
}
