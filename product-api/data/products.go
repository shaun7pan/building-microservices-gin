package data

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the structure of an API product
type Product struct {
	ID          int     `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float32 `json:"price" binding:"gt=0"`
	SKU         string  `json:"sku" binding:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

var Sku validator.Func = func(fl validator.FieldLevel) bool {
	// sku is of format abc-asds-sdfsdf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	maches := re.FindAllString(fl.Field().String(), -1)

	if len(maches) != 1 {
		return false
	}
	return true
}

var ErrProductNotFound = fmt.Errorf("Product not found")

// Products is a collection of product
type Products []*Product

// GetProducts return a list of products
func GetProducts() Products {
	return productList
}

func AddProducts(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProducts(id int, p *Product) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = p
	return nil
}

func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

func DeleteProductByID(id int) error {
	i := findIndexByProductID(id)

	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:(i-1)], productList[i:]...)
	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func getNextID() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd33",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
