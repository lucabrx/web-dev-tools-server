package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wdt/internal/data"
	"github.com/wdt/internal/tokens"
	validator "github.com/wdt/internal/validators"
)


func (app *application) registerUserWithMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	if !session.IsAnonymous() {
		app.alreadyHaveSessionResponse(w, r)
		return
	}

	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	user := &data.User{
		Email: input.Email,
		Name:  strings.Split(input.Email, "@")[0],
	}
	v := validator.New()
	data.ValidateEmail(v, input.Email)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	dbUser, err := app.models.Users.Get(0, user.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			err = app.models.Users.Insert(user)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		} else {
			app.serverErrorResponse(w, r, err)
		}
	} else {
		user = dbUser
	}

	magicLinkToken, err := tokens.CreateMagicLinkToken(user.Email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	magicLink := URL + "/v1/auth/magic-link/" + magicLinkToken
	emailData := struct {
		MagicLinkToken string
	}{
		MagicLinkToken: magicLink,
	}

	err = app.sendEmail("./templates/magic-email.tmpl", emailData, user.Email, "Magic Link - Web Dev Tools")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "success"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	
}

func (app *application) authenticateUserWithMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	if !session.IsAnonymous() {
		app.alreadyHaveSessionResponse(w, r)
		return
	}

	param := chi.URLParam(r, "token")

	email, err := tokens.ValidateMagicLinkToken(param)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.Get(0, email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := app.sessionCookie(token.Plaintext, token.Expiry)
	http.SetCookie(w, cookie)

	http.Redirect(w, r, app.config.ClientAddress, http.StatusFound)
}


func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := app.contextGetUser(r)

	err := app.models.Tokens.DeleteAllForUser(data.ScopeAuthentication, session.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := app.sessionCookie("", time.Unix(0, 0))
	http.SetCookie(w, cookie)

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "you have been logged out"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}