package handlers

import (
	"net/http"

	"github.com/AdiAkhileshSingh15/microservices-productapi/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//
//		200: productResponse
//	 422: errorValidation
//	 501: errorResponse

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Error("Handle POST Product")

	rw.Header().Add("Content-Type", "application/json")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Info("Prod: %#v", prod)
	data.AddProduct(prod)
}
