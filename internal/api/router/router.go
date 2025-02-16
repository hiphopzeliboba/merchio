package router

import (
	"github.com/go-chi/chi/v5"
	"merchio/internal/api/handler"
)

func NewRouter(userHandler *handler.Implementation) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/auth", userHandler.AuthHandler)
	//	r.Get("/api/sendCoin", middleware.AuthMiddleware(userHandler.SendCoinHandler))  //Unimplemented
	return r
}
