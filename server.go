package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int

	// user list
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// brodcast channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}
func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}
func (s *Server) Handler(connect net.Conn) {
	// user online!!
	user := NewUser(connect)
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// broadcast to other user
	s.Broadcast(user, "Im online")

	select {}
}

func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("error listen: ", err)
		return
	}
	defer listener.Close()

	go s.ListenMessage()
	// accept
	for {
		connect, err := listener.Accept()
		if err != nil {
			fmt.Println("error type: ", err)
			continue
		} else {
			fmt.Println("accept ", connect)
		}
		// handle
		go s.Handler(connect)
	}

}
