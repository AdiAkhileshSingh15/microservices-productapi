package handlers

import (
	"net/http"

	"github.com/AdiAkhileshSingh15/microservices-productapi/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Handle PUT Product")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	rw.Header().Add("Content-Type", "application/json")

	p.l.Info("Prod: %#v", prod)
	err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
