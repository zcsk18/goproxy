package main

import (
	"fmt"
	"goproxy/core"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)


func init() {
	f, _ := os.OpenFile("goproxy_clt" + ".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	log.SetOutput(f)
	log.SetOutput(os.Stdout)
}

func handleClient(c net.Conn) {
	defer c.Close()

	buff := make([]byte, 1024*10)
	length,err := c.Read(buff)
	if err != nil {
		return
	}

	method := core.GetMethod(buff[:length])
	host := core.GetHost(buff[:length])

	srv, err := net.DialTimeout("tcp", core.ProxyIP+":"+core.ProxyPort, 60*time.Second)
	if err != nil {
		log.Println("connect srv err")
		return
	}
	defer srv.Close()

	core.ProxySend(srv, "c", []byte(host))
	_, err = core.ProxyRecv(srv)
	if err != nil {
		log.Println("connect remote err")
	}

	if method == "CONNECT" {
		c.Write([]byte("HTTP/1.0 200 Connection established\r\n\r\n"))
		hijacker, ok := c.(http.Hijacker)
		if !ok {
			return
		}

		c, _, err = hijacker.Hijack()
		if err != nil {
			return
		}

	} else {
		core.ProxySend(srv, "p", buff[:length])
	}

	go core.ProxyTransfer(srv, c)
	for {
		msg, err := core.ProxyRecv(srv)
		if err != nil {
			return
		}
		switch msg.Op {
		case "p":
			c.Write(msg.Data)
		}
	}

}

func main() {
	port := core.LocalPort

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		log.Printf("new clt : %s \n", c.RemoteAddr());
		go handleClient(c)
	}
}