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

// AfterSocketEvent : this sets up a listener to a socket, and calls the play function only after the socket has received an event
// caller function can set the size of the expected event message - defaults to 512 if set to lesser than that
// returns a func to stop the event loop, and error incase any
func AfterSocketEvent(sock string, play func(msg []byte), size int) (func(), error) {
	if size < 512 {
		// incoming message size if not set by the caller code - it defaults to 512
		size = 512
	}
	start, stop, err := ListenOnUnixSocket(sock, func(c net.Conn) {
		buf := make([]byte, size)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		data := buf[0:nr]
		play(data)
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to setup AfterSocket event :%s", err)
	}
	go start()
	return stop, nil
}
