package utilities

import (
	"fmt"
	"net"
	"os"
)

// SendOverUnixSocket : starts a task that when receives a message on a channel can dispatch a message on socket
func SendOverUnixSocket(cancel chan interface{}, sock string, errx chan error) (chan []byte, func()) {
	send := make(chan []byte, 10) // client code to use this channel to send desired message
	return send, func() {
		defer close(send)
		for {
			select {
			case <-cancel:
				return
			case msg := <-send:
				conn, err := net.Dial("unix", sock)
				if err != nil {
					errx <- fmt.Errorf("Autolumin: Error connecting to srvrelay over socket %s", err)
				}
				_, err = conn.Write(msg)
				if err != nil {
					errx <- fmt.Errorf("Error writing message to TCP sock %s", err)
				}
				conn.Close()
			}
		}
	}
}

// ListenOnUnixSocket : sets up a listener, returns a function to start, stop and error incase the socket could not be listened
// https://eli.thegreenplace.net/2019/unix-domain-sockets-in-go/
func ListenOnUnixSocket(sock string, handler func(net.Conn)) (func(), func(), error) {
	os.RemoveAll(sock)
	l, err := net.Listen("unix", sock)
	if err != nil {
		return nil, nil, err
	}
	start := func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				return
			}
			go handler(fd)
		}
	}
	stop := func() {
		l.Close()
	}
	return start, stop, nil
}
