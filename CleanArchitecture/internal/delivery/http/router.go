package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter - HTTP 라우터 설정
func NewRouter(userHandler *UserHandler) *chi.Mux {
	r := chi.NewRouter()

	// 미들웨어
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// 헬스 체크
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API 라우트
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", userHandler.GetAllUsers)
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	return r
}

