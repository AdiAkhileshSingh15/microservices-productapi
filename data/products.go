package data

import (
	"context"
	"fmt"
	"time"

	protos "github.com/AdiAkhileshSingh15/microservices-currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for this product
	//
	// required: true
	// min: 1
	ID int `json:"id"`
	// the name for this product
	//
	// required: true
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	// the price for this product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"gt=0"`
	// the sku for this product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"required,sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// Products defines a slice of Product
type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductsDB(cc protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{cc, l, make(map[string]float64), nil}

	go pb.handleUpdates()

	return pb
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("Error subscribing to rates", "error", err)
		return
	}

	p.client = sub

	for {
		srr, err := sub.Recv()

		if err != nil {
			p.log.Error("Error while waiting for message", "error", err)
			return
		}

		// handle a returned error message
		if grpcError := srr.GetError(); grpcError != nil {
			sre := status.FromProto(grpcError)

			if sre.Code() == codes.InvalidArgument {
				errDetails := ""
				// get the RateRequest serialized in the error response
				// Details is a collection but we are only returning a single item
				if d := sre.Details(); len(d) > 0 {
					p.log.Error("Deets", "d", d)
					if rr, ok := d[0].(*protos.RateRequest); ok {
						errDetails = fmt.Sprintf("base: %s destination: %s", rr.GetBase().String(), rr.GetDestination().String())
					}
				}

				p.log.Error("Received error from currency service rate subscription", "error", grpcError.GetMessage(), "details", errDetails)
			}
		}

		if resp := srr.GetRateResponse(); resp != nil {
			p.log.Info("Received updated rate from server", "dest", resp.Destination.String(), "rate", resp.Rate)

			if err != nil {
				p.log.Error("Error receiving message", "error", err)
				return
			}

			p.rates[resp.Destination.String()] = resp.Rate
		}

	}
}

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// GetProducts returns all products from the database
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("getting rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}

	return pr, nil
}

func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}
	if currency == "" {
		return productList[i], nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("getting rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *productList[i]
	np.Price = np.Price * rate

	return &np, nil
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func UpdateProduct(p *Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = p

	return nil
}

func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	// if r, ok := p.rates[destination]; ok {
	// 	return r, nil
	// }

	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	// get initial rate
	resp, err := p.currency.GetRate(context.Background(), rr)
	if err != nil {
		// convert the GRPC error message
		grpcError, ok := status.FromError(err)
		if !ok {
			// unable to convert grpc error
			return -1, err
		}

		// if this is an Invalid Arguments exception santise the message before returning
		if grpcError.Code() == codes.InvalidArgument {
			return -1, fmt.Errorf("unable to retreive exchange rate from currency service: %s", grpcError.Message())
		}
	}

	p.rates[destination] = resp.Rate

	// subscribe for updates
	p.client.Send(rr)

	return resp.Rate, err
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
}
