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

// MongoConnectInfo: unified common interface that lets you connect to mongo irespective of how you choose the connection params
type MongoConnectInfo interface {
	Connect() (*mongo.Client, error) // ApplyURI to connect to mongo, user password combination
	URi() string
}

/* Mongo connection using the complete connection string */
type MongoConnectString string

func (mcs MongoConnectString) Connect() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return mongo.Connect(ctx, options.Client().ApplyURI(string(mcs)))
}
func (mcs MongoConnectString) URi() string {
	return string(mcs)
}

/* Mongo connection using mongo connection params */
type MongoConnectParams struct {
	Server string
	User   string
	Passwd string
}

func (mcp MongoConnectParams) URi() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", mcp.User, mcp.Passwd, mcp.Server)
}

func (mcp MongoConnectParams) Connect() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return mongo.Connect(ctx, options.Client().ApplyURI(mcp.URi()))
}

/* Use this as middleware, pass the connnect info and database name  */

func MongoConnect(mcinfo MongoConnectInfo, dbname string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cl, err := mcinfo.Connect()
		if err != nil {
			httperr.HttpErrOrOkDispatch(c, httperr.ErrGatewayConnect(err), log.WithFields(log.Fields{
				"stack": "MongoConnect",
				"uri":   mcinfo.URi(),
			}))
			return
		}
		c.Set("mongo-client", cl)
		c.Set("mongo-database", cl.Database(dbname))
	}
}

// MongoPingTest : will do a simple ping test to see if the server is reachable
func MongoPingTest(mcinfo MongoConnectInfo) error {
	cl, err := mcinfo.Connect()
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = cl.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to ping server %s", err)
	}
	log.Info("database is reachable..")
	defer cl.Disconnect(ctx) // purpose of this connection is served
	return nil
}
