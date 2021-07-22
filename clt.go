package main

import (
	"fmt"
	"goproxy/core"
	"net"
	"strconv"
)

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
	defer srv.Close()

	if err != nil {
		return
	}

	core.ProxySend(srv, "t", []byte(core.Token))
	_, err = core.ProxyRecv(srv)
	if err != nil {
		return
	}

	core.ProxySend(srv, "c", []byte(target))
	_, err = core.ProxyRecv(srv)
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
		buff := core.Pool.Get().([]byte)
		defer core.Pool.Put(buff)

		for {
			n, err := srv.Read(buff)
			if err != nil {
				return
			}
			core.Decode(&buff, n)
			client.Write(buff[:n])
		}
	}()

	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	for{
		n,err := client.Read(buff)
		if err != nil {
			return
		}
		core.Encode(&buff, n)
		srv.Write(buff[:n])
	}
}


