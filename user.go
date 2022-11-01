package main

import (
	"fmt"
	"net"
)

// 创建User类管理服务器的连接
type User struct {
	Name    string
	Addr    string
	C       chan string
	connect net.Conn
}

func NewUser(connect net.Conn) *User {
	userAddr := connect.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		C:       make(chan string),
		connect: connect,
	}
	go user.ListenMessage()
	return user
}
// 向客户端发送消息
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.connect.Write([]byte(msg + "\n"))
	}
}
