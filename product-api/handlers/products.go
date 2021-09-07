// Package classification of Product API
//
// Documetation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Products struct {
	l *log.Logger
}

type Param struct {
	ID int `uri:"id" form:"id" binding:"required,gt=0"`
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductIDFromURI(c *gin.Context) (int, error) {
	param := Param{}

	if err := c.ShouldBindUri(&param); err != nil {
		return -1, err
	}

	return param.ID, nil
}

func getProductIDFromParam(c *gin.Context) (int, error) {
	param := Param{}

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Println(err.Error())
		return -1, err
	}
	return param.ID, nil
}
