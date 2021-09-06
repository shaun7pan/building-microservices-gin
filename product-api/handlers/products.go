package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaun7pan/building-microservices-gin/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(c *gin.Context) {
	p.l.Println("Handle GET Product")

	// fetch the products from datastore
	ps := data.GetProducts()
	c.JSON(http.StatusOK, ps)
}

func (p *Products) AddProduct(c *gin.Context) {
	p.l.Printf("Handle POST product")

	// fetch new product from request
	newProd := &data.Product{}
	err := c.ShouldBindJSON(newProd)
	if err != nil {
		// http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	data.AddProducts(newProd)
}

func (p *Products) UpdateProducts(c *gin.Context) {
	p.l.Println("Handling PUT requests.")

	prod := &data.Product{}
	err := c.ShouldBindJSON(prod)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}

	// fetch id from URI
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.UpdateProducts(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(c.Writer, data.ErrProductNotFound.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
