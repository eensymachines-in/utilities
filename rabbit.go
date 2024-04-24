package utilities

/* Gin middleware to inject rabbitmq connection to gin context */
import (
	"fmt"

	"github.com/eensymachines-in/errx/httperr"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RabbitDirectXchnge : middleware handler can connect to rabbit mq and inject aqmp ch and cnnection in context
// server	: server hostname with port
// user pass : connecting credentials
// name of the xchange
func RabbitDirectXchnge(server, user, pass, xchangName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		connString := fmt.Sprintf("amqp://%s:%s@%s", user, pass, server)
		conn, err := amqp.Dial(connString)
		if err != nil { //failed to connect to gateway amqp
			httperr.HttpErrOrOkDispatch(c, httperr.ErrGatewayConnect(err), log.WithFields(log.Fields{
				"stack":       "RabbitConnectWithChn",
				"conn_string": connString,
			}))
			return
		}
		ch, err := conn.Channel()
		if err != nil { // failed to create channel
			httperr.HttpErrOrOkDispatch(c, httperr.ErrGatewayConnect(err), log.WithFields(log.Fields{
				"stack": "RabbitConnectWithChn",
			}))
			conn.Close() // incase no channel, we close the channel before we exit the stack
			return
		}
		err = ch.ExchangeDeclare(
			xchangName, // name
			"direct",   // exhange type
			true,       // durable
			false,      //auto deleted
			false,      //internal
			false,      // nowait
			nil,        //amqp.table
		)
		if err != nil { // failed to declare channel
			httperr.HttpErrOrOkDispatch(c, httperr.ErrGatewayConnect(err), log.WithFields(log.Fields{
				"stack": "RabbitConnectWithChn",
			}))
			// incase declaring the exchange fails we close the channel and connection on our way out
			ch.Close()
			conn.Close()
			return
		}
		c.Set("amqp-ch", ch)
		c.Set("amqp-conn", conn)
		c.Next()
	}
}
