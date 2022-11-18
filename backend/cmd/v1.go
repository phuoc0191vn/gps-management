package main

import (
	"net/http"

	"ctigroupjsc.com/phuocnn/gps-management/app/api"
	"github.com/julienschmidt/httprouter"
)

func NewAPIv1(container *Container) http.Handler {
	router := api.NewRouter()

	v1 := router.Group("/api/v1")

	authRouter(v1)

	v1.Use(
		api.RequireAuth(container.Config.JwtSecret),
		func(handle httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				handle(w, r, p)
			}
		},
	)

	userRouter(v1)
	accountRouter(v1)

	return router
}

func authRouter(parent *api.Router) {
	authHandler := api.AuthHandler{
		JwtSecret:          container.Config.JwtSecret,
		AccountRepository:  container.AccountRepository,
		AccessIDRepository: container.UserAccessIDRepository,
	}

	parent.POST("/auth/login", authHandler.BasicLogin)
}

func userRouter(parent *api.Router) {
	userHandler := api.UserHandler{
		UserRepository: container.UserRepository,
	}

	router := parent.Group("/user")

	router.GET("/:id/id", userHandler.GetByID)

	router.POST("", userHandler.AddUser)
}

func accountRouter(parent *api.Router) {
	accountHandler := api.AccountHandler{
		AccountRepository: container.AccountRepository,
	}

	router := parent.Group("/account")

	router.POST("", accountHandler.Add)
}
