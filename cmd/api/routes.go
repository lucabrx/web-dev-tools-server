package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.rateLimit)
	r.Use(app.enableCORS)
	r.Use(app.authenticate)

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Route("/v1/tools", func(r chi.Router) {
		r.Post("/", app.requireAuthenticatedUser(app.createToolHandler))
		r.Get("/{id}", app.requireAuthenticatedUser(app.getToolHandler))
		r.Delete("/{id}", app.adminPermission(app.requireAuthenticatedUser(app.deleteToolHandler)))
		r.Patch("/{id}", app.adminPermission(app.requireAuthenticatedUser(app.updateToolHandler)))
		r.Get("/", app.getToolsHandler)
		r.Get("/admin", app.adminPermission(app.requireAuthenticatedUser(app.getAdminToolsHandler)))
		r.Get("/toggle-published/{id}", app.adminPermission(app.requireAuthenticatedUser(app.toggleToolPublishedHandler)))
	})

	r.Route("/v1/categories", func(r chi.Router) {
		r.Post("/", app.requireAuthenticatedUser(app.createCategoryHandler)) 
		r.Get("/", app.getCategoriesHandler)
		r.Get("/admin", app.adminPermission(app.requireAuthenticatedUser(app.getAdminCategoriesHandler)))
		r.Delete("/{id}", app.adminPermission(app.requireAuthenticatedUser(app.deleteCategoryHandler)))
		r.Get("/toggle-published/{id}", app.adminPermission(app.requireAuthenticatedUser(app.toggleCategoryPublishedHandler)))
	})

	r.Route("/v1/upload", func(r chi.Router) {
		r.Post("/image", app.requireAuthenticatedUser(app.toggleCategoryPublishedHandler))
	})


	r.Get("/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]string{
			"status": "ok",
		}

		app.writeJSON(w, http.StatusOK, res, nil)
	})

	return r
}