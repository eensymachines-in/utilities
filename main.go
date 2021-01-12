package utilities

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// sysSigListener : sets up a listener to the system signals and when the signal actually occurs this will indicate by closing the cancel channel
// We intend to make a closure of this sort so that we can re-use the same
func sysSigListener() (func(), chan interface{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	interrupt := make(chan interface{}, 1) // system interrupt is communicated by closing this channel
	return func() {
		defer close(sigs)
		for {
			select {
			case <-sigs:
				log.Warn("System system interruption")
				log.Warn("Now closing all tasks..")
				close(interrupt)              //this will indicate to the main process that its time to close all the parent tasks
				<-time.After(1 * time.Second) //let all the tasks complete closing and cleanup
				return
			}
		}
	}, interrupt
}
