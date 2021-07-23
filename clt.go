package main

import (
	"fmt"
	"goproxy/cipher"
	"goproxy/core"
	"goproxy/proto"
	"goproxy/utils"
	"net"
)

var iniParser utils.IniParser

func init() {
	err := iniParser.Load("conf.ini")
	if err != nil {
		panic("cant find conf.ini")
	}
}

func main() {
	port := iniParser.GetString("clt", "port")
	if port == "" {
		panic("cant find conf port")
	}

	l, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic("listen err :" + port)
		return
	}

	pt, err := proto.GetDriver(iniParser)
	if err != nil {
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			panic("accept err")
			continue
		}
		go process(c, pt)
	}
}



func process(c net.Conn, pt proto.Driver) {
	defer c.Close()

	err := pt.Auth(c)
	if err != nil {
		return;
	}

	target, err := pt.GetTarget(c)
	if err != nil {
		return
	}

	cip , err := cipher.GetDriver(iniParser)
	if err != nil {
		return
	}

	s, err := core.Connect(iniParser.GetString("srv", "host"), iniParser.GetString("srv", "port"), cip)
	if err != nil {
		return
	}
	defer s.Close()

	_, err = s.Send([]byte(iniParser.GetString("common", "token")))
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






/*
func main() {
	server, err := net.Listen("tcp", ":"+strconv.Itoa(core.CltPort))
	if err != nil {
		fmt.Printf("Listen failed: %v\n", err)
		return
	}

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Printf("Accept failed: %v", err)
			continue
		}
		go process(client)
	}
}


func process(client net.Conn) {
	defer client.Close()

	if err := core.Socks5Auth(client); err != nil {
		fmt.Println("auth error:", err)
		return
	}

	target, err := core.Socks5Target(client)
	if err != nil {
		fmt.Println("get target error:", err)
		return
	}


	destAddrPort := fmt.Sprintf("%s:%d", core.SrvHost, core.SrvPort)
	srv, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return
	}
	defer srv.Close()

	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	_, err = srv.Write([]byte(core.Token))
	if err != nil {
		return
	}
	_, err = srv.Read(buff)
	if err != nil {
		return
	}

	_, err = srv.Write([]byte(target))
	if err != nil {
		return
	}

	_, err = srv.Read(buff)
	if err != nil {
		return
	}

	fmt.Printf("new connect from %s to %s \n", client.RemoteAddr(), target)
	_, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return
	}

	go func() {
		defer srv.Close()
		defer client.Close()

		for {
			buff, n, err := core.ProxyRead(srv)
			if err != nil {
				return
			}
			client.Write(buff[:n])
		}
	}()

	for{
		n,err := client.Read(buff)
		if err != nil {
			return
		}
		core.ProxyWrite(srv, buff[:n])
	}
}
*/

