package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/wdt/internal/data"
	"net/http"
	"strconv"
)

func (app *application) addFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Favorites.AddFavorite(session.ID, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"favorite": true}, nil)
}

func (app *application) removeFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)
	params := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Favorites.RemoveFavorite(session.ID, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"favorite": false}, nil)
}

func (app *application) getFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	favorites, err := app.models.Favorites.GetFavorites(session.ID)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var tools []*data.Tool

	for _, favorite := range favorites {
		tool, err := app.models.Tools.Get(favorite.ToolId)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		tool.Favorite = true
		tools = append(tools, tool)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tools": tools}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
