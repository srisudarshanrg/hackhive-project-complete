package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/config"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/carboprint", handlers.Repo.CarboPrint)

	mux.Get("/recycle-locator", handlers.Repo.RecycleLocator)

	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)

	mux.Get("/sign-up", handlers.Repo.SignUp)
	mux.Post("/sign-up", handlers.Repo.PostSignUp)

	mux.Get("/forgot-password", handlers.Repo.ForgotPassword)
	mux.Post("/forgot-password", handlers.Repo.PostForgotPassword)

	mux.Get("/otp-confirm", handlers.Repo.ConfirmOTP)
	mux.Post("/otp-confirm", handlers.Repo.PostConfirmOTP)

	mux.Get("/reset-password", handlers.Repo.ResetPassword)
	mux.Post("/reset-password", handlers.Repo.PostResetPassword)

	mux.Get("/resource-status", handlers.Repo.ResourceStatus)
	mux.Post("/resource-status", handlers.Repo.PostResourceStatus)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
