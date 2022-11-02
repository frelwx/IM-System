package main

import (
	"flag"
	"fmt"
	"net"
)
type Client struct {
	ServerIP string
	ServerPort int
	Name string
	connect net.Conn
}

func NewClient(ServerIP string, ServerPort int) *Client {
	connect, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIP, ServerPort))
	if(err != nil) {
		fmt.Println("dial error", err)
		return nil
	}
	client := &Client{
		ServerIP: ServerIP,
		ServerPort: ServerPort,
		Name: "tmp",
		connect: connect,
	}
	return client
}

var serverIp string
var serverPort int
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set ip(default is 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 9190, "set port(default is 9190)")
	flag.Parse()
}
func main() {
	client := NewClient(serverIp, serverPort)
	if(client == nil) {
		fmt.Println(">>>>>>connect error")
		return
	}
	fmt.Println(">>>>>connect success!")
	select{}
}