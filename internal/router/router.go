package router

import (
	"Social-Net-Dialogs/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewRouter(i *handler.Instance) http.Handler {
	r := chi.NewRouter()
	r.Mount("/debug", middleware.Profiler())
	r.Mount("/metrics", promhttp.Handler())
	r.Group(func(r chi.Router) {
		r.Use(i.BasicAuth)

		r.Get("/dialog/{user_id}/list", i.GetMessages)
		r.Post("/dialog/{user_id}/send", i.SendMessage)
		r.Put("/dialog/{user_id}/list/{messageId}", i.MarkAsRead)
	})

	return r
}
