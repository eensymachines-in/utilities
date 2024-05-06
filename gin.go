package utilities

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/* restrictive  origin policy, you can set only on domain, with port or aubdomain that you want to restrict the requests from */
func CorsWithOrigin(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, we add the headers with need to enable CORS
		// Make sure to adjust these headers to your needs
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Content-Type", "application/json")
		// Second, we handle the OPTIONS problem
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			// Everytime we receive an OPTIONS request,
			// we just return an HTTP 200 Status Code
			// Like this, Angular can now do the real
			// request using any other method than OPTIONS
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

/* Wide open cors policy, no restrictions on the origin, but all origins have to be working on port 80 */
func CORS(c *gin.Context) {
	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		// Everytime we receive an OPTIONS request,
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}

// CORS : this allows all cross origin requests
