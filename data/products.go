package data

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

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
	Price float32 `json:"price" validate:"gt=0"`
	// the sku for this product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"required,sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

type Products []*Product

var ErrProductNotFound = fmt.Errorf("Product not found")

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	sku := fl.Field().String()
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(sku, -1)
	if len(matches) != 1 {
		return false
	}

	return true
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

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
