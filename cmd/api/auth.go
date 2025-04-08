package main

import (
	"github.com/seanhalberthal/webmart/internal/store"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// RegisterUser godoc
//
//	@Summary	Registers a new user on the system
//	@Tags		authentication
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		RegisterUserPayload	true	"User credentials"
//	@Success	201		{object}	store.User			"User registered"
//	@Failure	400		{object}	error
//	@Failure	500		{object}	error
//	@Security	ApiKeyAuth
//	@Router		/users [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	user := &store.User{Username: payload.Username, Email: payload.Email}

	if err := user.Password.Set(payload.Password); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	ctx := r.Context()

	err := app.store.Users.UserCreate(ctx, user)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	if err := writeJSONResponse(w, http.StatusCreated, payload); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}
