package main

import (
	"github.com/google/uuid"
	"github.com/seanhalberthal/webmart/internal/store"
	"net/http"
)

type CreateListingPayload struct {
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rating      int       `json:"rating"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
}

func (app *application) createListingHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateListingPayload
	if err := readJSON(w, r, &payload); err != nil {
		err := respondWithErrorJSON(w, http.StatusBadRequest, err.Error())
		if err != nil {
			return
		}
		return
	}

	listing := &store.Product{
		UserID:      payload.UserID,
		Title:       payload.Title,
		Description: payload.Description,
		Rating:      payload.Rating,
		Price:       payload.Price,
		Stock:       payload.Stock,
	}

	ctx := r.Context()

	if err := app.store.Products.ProductCreate(ctx, listing); err != nil {
		err := respondWithErrorJSON(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, listing); err != nil {
		err := respondWithErrorJSON(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
}
