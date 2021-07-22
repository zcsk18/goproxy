package main

import (
	"fmt"
	"goproxy/core"
	"log"
	"net"
	"os"
	"time"
)

func init() {
	f, _ := os.OpenFile("goproxy_srv" + ".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	log.SetOutput(f)
	log.SetOutput(os.Stdout)
}

func handleConn(c net.Conn) {
	defer c.Close()
	var dest net.Conn

	for {
		msg, err := core.ProxyRecv(c)
		if err != nil {
			break;
		}
		switch msg.Op {
		case "c" :
			dest, err = net.DialTimeout("tcp", string(msg.Data), 60*time.Second)
			if err != nil {
				log.Println("connect remote err")
				return
			}
			log.Printf("conenct %s suc \n", msg.Data)
			core.ProxySend(c, "o", []byte("ok"))
			go core.ProxyTransfer(c, dest)
		case "p" :
			dest.Write([]byte(msg.Data))
		}
	}

	log.Printf("connection closed : %s \n", c.RemoteAddr())
}


func main() {
	port := core.ProxyPort

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
		go handleConn(c)
	}

}

