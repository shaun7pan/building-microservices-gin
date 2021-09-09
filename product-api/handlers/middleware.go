package handlers

import (
	"github.com/gin-gonic/gin"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// Custom middleware
func (p *Products) CustomMiddleware(c *gin.Context) {
	p.l.Println("Custom middleware log.")
	c.Next()
}

// create custom middleware
func (p *Products) BuildCustomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		p.l.Println("Username in header is: ", c.Request.Header.Get("username"))
		c.Next()
	}
}
