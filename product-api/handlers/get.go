package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaun7pan/building-microservices-gin/product-api/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all products
func (p *Products) ListAll(c *gin.Context) {
	p.l.Println("Handle GET Products")

	// fetch the products from datastore
	ps := data.GetProducts()
	c.JSON(http.StatusOK, ps)
}

// swagger:route GET /products/{id} products listSingle
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse
//	500: internalErrorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(c *gin.Context) {
	id, err := getProductIDFromURI(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	p.l.Println("[DEBUG] get product id", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)
		c.JSON(http.StatusNotFound, gin.H{
			"Message": err.Error(),
		})
		return
	default:
		p.l.Println("[ERROR] fetching product", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, prod)
}
