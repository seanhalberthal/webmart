package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/seanhalberthal/webmart/internal/store"
	"net/http"
)

type CreateProductPayload struct {
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rating      int       `json:"rating"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
}

func getProductID(w http.ResponseWriter, r *http.Request) uuid.UUID {
	idStr := chi.URLParam(r, "productID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(w, http.StatusBadRequest, err)
	}
	return id
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateProductPayload
	if err := readJSON(w, r, &payload); err != nil {
		handleError(w, http.StatusBadRequest, err)
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
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, listing); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := getProductID(w, r)

	product, err := app.store.Products.ProductGet(ctx, id)
	if err != nil {
		handleError(w, http.StatusNotFound, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := getProductID(w, r)

	err := app.store.Products.ProductDelete(ctx, id)
	if err != nil {
		handleError(w, http.StatusNotFound, err)
	}

	deletedProduct := map[string]string{
		"message":   "Product deleted",
		"productID": id.String(),
	}

	if err := writeJSONResponse(w, http.StatusOK, deletedProduct); err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}
}
