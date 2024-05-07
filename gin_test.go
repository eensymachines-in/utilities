package utilities

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOriginPatterns(t *testing.T) {
	patt := regexp.MustCompile(`^http://[a-zA-Z.]*eensymachines.in[:0-9]*[\/]*$`)

	yes := patt.MatchString("http://aqua.eensymachines.in")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://aqua.eensymachines.in/")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://eensymachines.in/")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://eensymachines.in")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://aqua.eensymachines.in:8080")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://aqua.eensymachines.in:8080/")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://eensymachines.in:8080/")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	patt = regexp.MustCompile(`^http://localhost[:0-9]*[\/]*$`)

	yes = patt.MatchString("http://localhost")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://localhost/")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://localhost:8080/")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")

	yes = patt.MatchString("http://localhost:8080")
	assert.Equal(t, yes, true, "Unexpected mismatch on the pattern")
}
