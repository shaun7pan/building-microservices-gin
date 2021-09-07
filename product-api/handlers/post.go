package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaun7pan/building-microservices-gin/product-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	202: productResponse
//  422: errorValidation
//  501: errorResponse

func (p *Products) Create(c *gin.Context) {
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

	c.Writer.WriteHeader(http.StatusNoContent)
}
