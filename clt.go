package main

import (
	"fmt"
	"goproxy/proto"
	"goproxy/utils"
	"net"
)

func main() {
	port := utils.Ini.GetString("clt", "port")
	if port == "" {
		panic("cant find conf port")
	}

	l, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic("listen err :" + port)
		return
	}

	pt, err := proto.GetDriver(utils.Ini)
	if err != nil {
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			panic("accept err")
			continue
		}
		fmt.Printf("new accept:%s\n", c.RemoteAddr())
		go pt.Process(c)
	}
}
