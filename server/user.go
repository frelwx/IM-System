package main

import (
	"fmt"
	"net"
	"strings"
)

// 创建User类管理服务器的连接
type User struct {
	Name    string
	Addr    string
	C       chan string
	connect net.Conn
	server  *Server
}

func NewUser(connect net.Conn, server *Server) *User {
	userAddr := connect.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		C:       make(chan string),
		connect: connect,
		server:  server,
	}
	go user.ListenMessage()
	return user
}

// 向客户端发送消息
func (u *User) ListenMessage() {
	for {
		msg, ok := <-u.C
		if !ok {
			fmt.Println("channel C closed")
			break
		}
		u.connect.Write([]byte(msg + "\n"))
	}
}

func (u *User) OnLine() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// broadcast to other user
	u.server.Broadcast(u, "Im online")
}
func (u *User) OffLine() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	close(u.C)
	// broadcast to other user
	u.server.Broadcast(u, "Im offline")
}

func (u *User) SendMsg(msg string) {
	u.connect.Write([]byte(msg))
}
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			who := "[" + user.Addr + "]" + user.Name + ": online...\n"
			u.SendMsg(who)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		NewName := msg[7:]
		fmt.Println(len(NewName), NewName)
		_, ok := u.server.OnlineMap[NewName]
		if ok {
			u.SendMsg("name[" + NewName + "] already used\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[NewName] = u
			u.server.mapLock.Unlock()
			u.Name = NewName
			u.SendMsg("sucessfully update name to" + NewName + "\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		msgSplit := strings.Split(msg, "|")
		ToUser, ok := u.server.OnlineMap[msgSplit[1]]
		if !ok {
			u.SendMsg(msgSplit[1] + "does not exist")
		}
		content := msgSplit[2]
		ToUser.SendMsg(u.Name + " says to you : " + content)

	} else {
		u.server.Broadcast(u, msg)
	}

}
