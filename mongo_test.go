package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongoPing(t *testing.T) {
	data := []MongoConnectInfo{
		MongoConnectString("mongodb://eensyaquap-dev:33n5y+4dm1n@aqua.eensymachines.in:32701"),
		MongoConnectParams{
			Server: "aqua.eensymachines.in:32701",
			User:   "eensyaquap-dev",
			Passwd: "33n5y+4dm1n",
		},
	}
	for _, d := range data {
		cl, err := d.Connect()
		assert.Nil(t, err, "Unexpected error when connecting to database")
		assert.NotNil(t, cl, "Unexpected nil client when connecting to databse")
		t.Log(d.URi())
		err = MongoPingTest(d)
		assert.Nil(t, err, "Unexpected error when ping testing the database")
	}
}
