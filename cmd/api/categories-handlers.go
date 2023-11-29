package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wdt/internal/data"
	validator "github.com/wdt/internal/validators"
)


func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.Category{
		Name: input.Name,
	}

	v := validator.New()
	data.ValidateCategories(v, category)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}


	err = app.models.Categories.Insert(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}


	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Categories.Delete(id)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "category successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) toggleCategoryPublishedHandler(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	category.Published = !category.Published

	err = app.models.Categories.Update(category)
	if err != nil {
		switch err {
		case data.ErrEditConflict:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := app.models.Categories.GetAllPublished()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAdminCategoriesHandler(w http.ResponseWriter, r *http.Request) {


	var meta struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	meta.Filters.Sort = app.readString(qs, "sort", "-id")
	meta.Filters.SortSafelist = []string{"name", "id", "-name", "-id", "published", "-published"}
	meta.Filters.Page = app.readInt(qs, "page", 1, v)
	meta.Filters.PageSize = app.readInt(qs, "pageSize", 20, v)

	if data.ValidateFilters(v, meta.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	categories, metadata, err := app.models.Categories.GetAll(meta.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories, "metadata" : metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
