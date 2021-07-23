package main

import (
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

	driver, err := proto.GetDriver(utils.Ini)
	if err != nil {
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			panic("accept err")
			continue
		}
		go driver.Process(c)
	}
}
