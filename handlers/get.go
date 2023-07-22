package handlers

import (
	"net/http"

	"github.com/AdiAkhileshSingh15/microservices-productapi/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//	200: productsResponse

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := data.ToJSON(lp, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
