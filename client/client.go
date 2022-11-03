package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	connect    net.Conn
	flag       int
}

func (c *Client) DealResponse() {
	io.Copy(os.Stdout, c.connect)
}
func (c *Client) menu() bool {
	fmt.Println("1.public chat")
	fmt.Println("2.private chat")
	fmt.Println("3.change name")
	fmt.Println("0.exit")

	var flag int
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>please input right number<<<<")
		return false
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {

		}

		switch c.flag {
		case 1:
			c.PublicChat()
		case 2:
			c.PrivateChat()
		case 3:
			c.UpdateName()
		case 0:
		}
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println("please input name")
	fmt.Scanln(&c.Name)

	SendMsg := "rename|" + c.Name + "\n"
	_, err := c.connect.Write([]byte(SendMsg))
	if err != nil {
		fmt.Println("connect write error", err)
		return false
	}
	return true
}
func (c *Client) PublicChat() {
	var SendMsg string
	fmt.Println("please input content, q for exit")
	fmt.Scanln(&SendMsg)
	for SendMsg != "q" {
		if len(SendMsg) > 0 {
			SendMsg = SendMsg + "\n"
			_, err := c.connect.Write([]byte(SendMsg))
			if err != nil {
				fmt.Println("public chat err", err)
				break
			}
		}
		fmt.Println("please input content, q for exit")
		fmt.Scanln(&SendMsg)
	}
}
func (c *Client) PrivateChat() {
	var target string
	var SendMsg string
	c.Who()
	fmt.Println("please input chat target, q for exit")
	fmt.Scanln(&target)
	
	for target != "q" {
		
		fmt.Println("please input content, q for exit")
		fmt.Scanln(&SendMsg)

		for SendMsg != "q" {
			if len(SendMsg) > 0 {
				SendMsg = "to|" + target + "|" + SendMsg + "\n\n"
				_, err := c.connect.Write([]byte(SendMsg))
				if err != nil {
					fmt.Println("public chat err", err)
					break
				}
			}
			fmt.Println("please input content, q for exit")
			fmt.Scanln(&SendMsg)
		}

		c.Who()
		fmt.Println("please input chat target, q for exit")
		fmt.Scanln(&target)
	}
}
func (c *Client) Who() {
	SendMsg := "who\n"
	_, err := c.connect.Write([]byte(SendMsg))
	if err != nil {
		fmt.Println("public chat err", err)
		return
	}
}
func NewClient(ServerIP string, ServerPort int) *Client {
	connect, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIP, ServerPort))
	if err != nil {
		fmt.Println("dial error", err)
		return nil
	}
	client := &Client{
		ServerIP:   ServerIP,
		ServerPort: ServerPort,
		Name:       "tmp",
		connect:    connect,
		flag:       4,
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
	if client == nil {
		fmt.Println(">>>>>>connect error")
		return
	}
	fmt.Println(">>>>>connect success!")
	go client.DealResponse()
	client.Run()
}
