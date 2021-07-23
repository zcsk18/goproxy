package main

import (
	"fmt"
	"goproxy/cipher"
	"goproxy/core"
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
		go process_clt(c, pt)
	}
}

func process_clt(c net.Conn, pt proto.Driver) {
	defer c.Close()

	err := pt.Auth(c)
	if err != nil {
		fmt.Printf("auth err: %s\n", err)
		return;
	}

	target, err := pt.GetTarget(c)
	if err != nil {
		return
	}

	cip , err := cipher.GetDriver(utils.Ini)
	if err != nil {
		fmt.Printf("cipher err: %s\n", err)
		return
	}

	fmt.Printf("srv %s:%s \n", utils.Ini.GetString("srv", "host"), utils.Ini.GetString("srv", "port"))
	s, err := core.Connect(utils.Ini.GetString("srv", "host"), utils.Ini.GetString("srv", "port"), cip)
	if err != nil {
		fmt.Printf("Connect err: %s\n", err)
		return
	}
	defer s.Close()

	_, err = s.Send([]byte(utils.Ini.GetString("common", "token")))
	if err != nil {
		return
	}

	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	_, err = s.Recv(buff)
	if err != nil {
		return
	}

	_, err = s.Send([]byte(target))
	if err != nil {
		return
	}

	_, err = s.Recv(buff)
	if err != nil {
		return
	}

	fmt.Printf("new connect from %s to %s \n", c.RemoteAddr(), target)
	err = pt.Connected(c)
	if err != nil {
		return
	}


	go func() {
		defer s.Close()
		defer c.Close()

		buff := core.Pool.Get().([]byte)
		defer core.Pool.Put(buff)

		for {
			n, err := s.Recv(buff)
			if err != nil {
				return
			}
			c.Write(buff[:n])
		}
	}()

	for{
		n,err := c.Read(buff)
		if err != nil {
			return
		}
		s.Send(buff[:n])
	}
}
