package routing

import (
	"log"
	"net/http"

	"github.com/egor_lukyanovich/avito/pkg/app"
	"github.com/go-chi/chi/v5"
)

func NewRouter(storage app.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Println("write response failed:", err)
		}

	})

	return r
}
