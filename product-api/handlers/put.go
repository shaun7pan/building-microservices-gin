package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaun7pan/building-microservices-gin/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) Update(c *gin.Context) {
	p.l.Println("Handling PUT requests.")

	prod := &data.Product{}
	err := c.ShouldBindJSON(prod)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}

	// fetch id from URI
	id, err := getProductIDFromURI(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
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

	c.Writer.WriteHeader(http.StatusNoContent)

}
