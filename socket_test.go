//go:build exclude

package utilities

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel) // default is info level, if verbose then trace
}

/*
handler: new connections are handled here, this will read the connection and print the message received
lets assume we receive json messages on this
*/
func handler(c net.Conn) {
	log.Printf("Client connected [%s]", c.RemoteAddr().Network())
	buff := make([]byte, 512)
	n, err := c.Read(buff)
	if err != nil {
		log.Print("Error copying data into buffer")
		log.Error(err)
		return
	}
	log.Printf("We have received about %d bytes of data", n)
	message := struct {
		Interrupt bool `json:"interrupt"`
	}{}
	err = json.Unmarshal(buff[0:n], &message)
	if err != nil {
		log.Print("failed to unmarshal message on the connection")
		log.Error(err)
		return
	}
	log.Printf("%v", message)
	c.Close()
}
func TestSockListener(t *testing.T) {
	start, stop, err := ListenOnUnixSocket("./test.sock", handler)
	go start()
	defer stop()
	c, err := net.Dial("unix", "./test.sock")
	if err != nil {
		t.Error(err)
		return
	}
	data, _ := json.Marshal(map[string]bool{"interrupt": true})
	c.Write(data)
	<-time.After(5 * time.Second)
}

func TestAfterSockEvent(t *testing.T) {
	stop, err := AfterSocketEvent("./test.sock", func(data []byte) {
		t.Log("We have received data over the test socket")
		t.Log(string(data))
	})
	defer stop()
	assert.Nil(t, err, "Unexpected error in setting up the after socket event")
	if err != nil {
		panic(err)
	}
	<-time.After(5 * time.Minute) //wait for netcat to send the data
}
