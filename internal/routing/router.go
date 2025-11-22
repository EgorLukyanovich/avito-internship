package routing

import (
	"net/http"

	"github.com/egor_lukyanovich/avito/pkg/app"
	"github.com/go-chi/chi/v5"
)

func NewRouter(storage app.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return r
}
