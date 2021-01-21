package utilities

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

// ListenOnUnixSocket : sets up a listener, returns a function to start, stop and error incase the socket could not be listened
// https://eli.thegreenplace.net/2019/unix-domain-sockets-in-go/
func ListenOnUnixSocket(sock string, handler func(net.Conn)) (func(), func(), error) {
	retries := 2
	var l net.Listener
	var err error
	for i := retries; i > 0; i-- {
		l, err = net.Listen("unix", sock)
		if err != nil {
			if i == 1 {
				return nil, nil, err
			}
			continue
		}
		break
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
		err := os.RemoveAll(sock) // for the next run we need the file to be removed .. this is part of server cleanup
		if err != nil {
			log.Errorf("Failed to clear socket %s", err)
		}
	}
	return start, stop, nil
}
