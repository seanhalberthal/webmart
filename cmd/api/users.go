package main

import (
	"github.com/seanhalberthal/webmart/internal/store"
	"net/http"
	"time"
)

type CreateUserPayload struct {
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		err := respondWithErrorJSON(w, http.StatusBadRequest, err.Error())
		if err != nil {
			return
		}
		return
	}

	user := &store.User{
		Name:      payload.Name,
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  payload.Password,
		CreatedAt: payload.CreatedAt,
	}

	ctx := r.Context()

	if err := app.store.Users.UserCreate(ctx, user); err != nil {
		err := respondWithErrorJSON(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, user); err != nil {
		err := respondWithErrorJSON(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
}
