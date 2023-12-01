package main

import (
	"errors"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wdt/internal/data"
	validator "github.com/wdt/internal/validators"
)

func (app *application) createToolHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
		ImageUrl    string `json:"imageUrl"`
		Website     string `json:"website"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tool := &data.Tool{
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
		ImageUrl:    input.ImageUrl,
		Published:   false,
		Website:     input.Website,
	}

	v := validator.New()
	if data.ValidateTools(v, tool); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tools.Insert(tool)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"tool": tool}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getToolHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tool, err := app.models.Tools.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tool": tool}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateToolHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Category    *string `json:"category"`
		Description *string `json:"description"`
		ImageUrl    *string `json:"imageUrl"`
		Published   *bool   `json:"published"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tool, err := app.models.Tools.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.Name != nil {
		tool.Name = *input.Name
	}
	if input.Category != nil {
		tool.Category = *input.Category
	}
	if input.Description != nil {
		tool.Description = *input.Description
	}
	if input.ImageUrl != nil {
		tool.ImageUrl = *input.ImageUrl
	}
	if input.Published != nil {
		tool.Published = *input.Published
	}

	v := validator.New()
	if data.ValidateTools(v, tool); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tools.Update(tool)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			v.AddError("version", "is not valid")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tool": tool}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteToolHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Tools.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) toggleToolPublishedHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tool, err := app.models.Tools.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	tool.Published = !tool.Published

	err = app.models.Tools.Update(tool)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tool": tool}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getToolsHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	var meta struct {
		data.Filters
		Search string
	}

	v := validator.New()
	qs := r.URL.Query()

	meta.Filters.Sort = app.readString(qs, "sort", "-id")
	meta.Filters.SortSafelist = []string{"name", "id", "-name", "-id", "published", "-published", "category", "-category"}
	meta.Filters.Page = app.readInt(qs, "page", 1, v)
	meta.Filters.PageSize = app.readInt(qs, "pageSize", 20, v)
	meta.Search = app.readString(qs, "search", "")

	tools, metadata, err := app.models.Tools.GetAllPublished(meta.Search, meta.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if session != nil {
		favorites, err := app.models.Favorites.GetFavorites(session.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		for _, tool := range tools {
			for _, favorite := range favorites {
				if tool.ID == favorite.ToolId {
					tool.Favorite = true
				}
			}
		}

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tools": tools, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAdminToolsHandler(w http.ResponseWriter, r *http.Request) {

	var meta struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	meta.Filters.Sort = app.readString(qs, "sort", "-id")
	meta.Filters.SortSafelist = []string{"name", "id", "-name", "-id", "published", "-published", "category", "-category"}
	meta.Filters.Page = app.readInt(qs, "page", 1, v)
	meta.Filters.PageSize = app.readInt(qs, "pageSize", 20, v)
	search := app.readString(qs, "search", "")

	tools, metadata, err := app.models.Tools.GetAll(meta.Filters, search)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if data.ValidateFilters(v, meta.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tools": tools, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
