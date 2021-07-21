package main

import (
	"goproxy/core"
	"log"
	"net"
	"os"
)

func init() {
	f, _ := os.OpenFile("goproxy_test" + ".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	log.SetOutput(f)
	log.SetOutput(os.Stdout)
}



func main() {
	str := "127.0.0.1:" + core.LocalPort
	c, err := net.Dial("tcp", string(str));
	if err != nil {
		log.Println("connect err")
		return;
	}

	s := "zcs"
	core.ProxySend(c, "c", []byte(s))
	msg, err := core.ProxyRecv(c)
	log.Println("recv: %s \n", msg.Data)
}