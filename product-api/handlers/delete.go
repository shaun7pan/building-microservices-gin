package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaun7pan/building-microservices-gin/product-api/data"
)

// swagger:route DELETE /products products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

func (p *Products) Delete(c *gin.Context) {

	//fetch id from uri
	id, err := getProductIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	p.l.Println("Deleting Product", "id", id)

	err = data.DeleteProductByID(id)
	if err != nil {
		p.l.Println("[Error] deleting product", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}

	if err == data.ErrProductNotFound {
		p.l.Println("[Error] deleting product", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"Message": err.Error(),
		})
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
