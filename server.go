package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}
	return server
}

func (s *Server) Handler(connect net.Conn) {
	fmt.Println("handling.......")
}

func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("error listen: ", err)
		return
	}
	defer listener.Close()
	// accept

	for {
		connect, err := listener.Accept()
		if err != nil {
			fmt.Println("error type: ", err)
			continue
		} else {
			fmt.Println("accept ", connect)
		}
		go s.Handler(connect)
	}
	// handle

	// close
}
