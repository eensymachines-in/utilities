package utilities

/* Defines go gin middleware functions that can inject mongo connection to downstream handlers
Use these to have mongo client and database as injected objects in gin context */
import (
	"context"
	"fmt"
	"time"

	"github.com/eensymachines-in/errx/httperr"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// USes connection string  and db name to connect
func MongoConnectURI(uri, dbname string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil || client == nil {
			httperr.HttpErrOrOkDispatch(c, httperr.ErrGatewayConnect(err), log.WithFields(log.Fields{
				"stack": "MongoConnect",
				"uri":   uri,
			}))
			return
		}
		c.Set("mongo-client", client)
		c.Set("mongo-database", client.Database(dbname))
	}
}

// Uses server, user password,dbname to connect to the mongo client
func MongoConnect(server, user, passwd, dbname string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", user, passwd, server)))
		if err != nil {
			httperr.HttpErrOrOkDispatch(c, httperr.ErrGatewayConnect(err), log.WithFields(log.Fields{
				"stack":  "MongoConnect",
				"login":  user,
				"server": server,
			}))
			return
		}
		if client.Ping(ctx, readpref.Primary()) != nil {
			httperr.HttpErrOrOkDispatch(c, httperr.ErrGatewayConnect(err), log.WithFields(log.Fields{
				"stack":  "MongoConnect",
				"login":  user,
				"server": server,
			}))
			return
		}
		c.Set("mongo-client", client)
		c.Set("mongo-database", client.Database(dbname))
	}
}

// MongoPingTest : will do a simple ping test to see if the server is reachable
func MongoPingTest(server, user, passwd string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", user, passwd, server)))
	if err != nil || client == nil {
		return fmt.Errorf("failed to connect to mongo server %s", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to ping server %s", err)
	}
	log.Info("database is reachable..")
	defer client.Disconnect(ctx) // purpose of this connection is served
	return nil
}
