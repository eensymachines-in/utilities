package utilities

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var (
	allowedOrigins = []*regexp.Regexp{
		regexp.MustCompile(`^http://[a-zA-Z.]*eensymachines.in[:0-9]*[\/]*$`),
		regexp.MustCompile(`^http://localhost[:0-9]*[\/]*$`),
		/* Add patterns of more origins that you would want to allow */
	}
)

/* restrictive  origin policy, you can set only on domain, with port or aubdomain that you want to restrict the requests from */

func CorsWithOrigin(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, we add the headers with need to enable CORS
		// Make sure to adjust these headers to your needs
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, PUT")
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

// populateHeaders : access-control-allow-origin is populated from the pattern
// methods allowed are set as header from one place - hence easier to maintain
// Allow headers also set in one place  - hence easier to matain
func populateHeaders(c *gin.Context) {
	/* Instead of making a wide open policy which is not only dangerous it fails to cover localhost in with different port
	xample localhost:8080 which is operating from a different port */
	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	reqHdrOrigin := c.Request.Header.Get("Origin")
	c.Header("Access-Control-Allow-Origin", "http://eensymachines.in") // this is the default if nothing is set in the for loop bwlo
	for _, patt := range allowedOrigins {
		if patt.MatchString(reqHdrOrigin) {
			c.Header("Access-Control-Allow-Origin", reqHdrOrigin)
			break
		}
	}
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, PUT")
	c.Header("Access-Control-Allow-Headers", "*")
}

func Preflight(c *gin.Context) {
	populateHeaders(c)
	c.Header("Access-Control-Max-Age", "86400")
	c.AbortWithStatus(http.StatusNoContent) //204
}

/* Wide open cors policy, no restrictions on the origin, but all origins have to be working on port 80 */
func CORS(c *gin.Context) {
	populateHeaders(c)
	c.Header("Content-Type", "application/json")
	// Second, we handle the OPTIONS problem
	c.Next()
}

// CORS : this allows all cross origin requests
