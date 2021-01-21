package utilities

import (
	"net"
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
		// the socket is removed in the client code not here.
		// we tried doing that but does not help, since the stop function is not allowed too much time
	}
	return start, stop, nil
}
