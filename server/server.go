package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(connect, s)
	user.OnLine()

	IsAlive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := connect.Read(buf)
			if n == 0 {
				fmt.Println("n == 0", err)
				user.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println(fmt.Println("conn message error", err))
				return
			}

			msg := string(buf[:n-1])
			fmt.Println("server received: ", msg)
			user.DoMessage(msg)
			IsAlive <- true
		}
	}()

	for {
		select {
		case <-IsAlive:
		case <-time.After(time.Second * 60 * 2):
			user.SendMsg("you are kicked!!")
			connect.Close()
			return
		}
	}

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
