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

// CreateProduct godoc
//
//	@Summary	Create a new product listing
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		body	body		CreateProductPayload	true	"Product creation payload"
//	@Success	201		{object}	store.Product
//	@Failure	400		{object}	error
//	@Failure	500		{object}	error
//	@Router		/products [post]
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

	if err := writeJSONResponse(w, http.StatusCreated, listing); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}

// GetProduct godoc
//
//	@Summary	Fetches a product by ID
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Product ID"
//	@Success	200	{object}	store.Product
//	@Failure	404	{object}	error
//	@Failure	500	{object}	error
//	@Router		/products/{id} [get]
func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := getProductID(w, r)

	product, err := app.store.Products.ProductGet(ctx, id)
	if err != nil {
		handleError(w, http.StatusNotFound, err)
		return
	}

	reviews, err := app.store.Reviews.ReviewGet(ctx, id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	product.Reviews = reviews

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}

// DeleteProduct godoc
//
//	@Summary	Delete product by ID
//	@Tags		products
//	@Param		id	path		int		true	"Product ID"
//	@Success	204	{string}	string	"Product deleted successfully"
//	@Failure	404	{object}	error
//	@Router		/products/{id} [delete]
func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := getProductID(w, r)

	err := app.store.Products.ProductDelete(ctx, id)
	if err != nil {
		handleError(w, http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateProductPayload struct {
	Title       string  `json:"title" validate:"omitempty,min=1,max=100"`
	Description string  `json:"description" validate:"omitempty,min=1,max=100"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
	Stock       int     `json:"stock" validate:"omitempty,min=0"`
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	Updates an existing product by ID
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Product ID"
//	@Param			body	body		UpdateProductPayload	true	"Updated product data"
//	@Success		200		{object}	store.Product
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/products/{id} [patch]
func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := getProductID(w, r)

	var payload UpdateProductPayload
	if err := readJSON(w, r, &payload); err != nil {
		handleError(w, http.StatusBadRequest, err)
		return
	}

	product, err := app.store.Products.ProductGet(ctx, id)
	if err != nil {
		handleError(w, http.StatusNotFound, err)
		return
	}

	product.Title = payload.Title
	product.Description = payload.Description
	product.Price = payload.Price
	product.Stock = payload.Stock

	err = app.store.Products.ProductUpdate(ctx, product)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	if err := writeJSONResponse(w, http.StatusOK, product); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}
