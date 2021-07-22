package main

import (
	"errors"
	"fmt"
	"goproxy/core"
	"net"
	"strconv"
)

func main() {
	server, err := net.Listen("tcp", ":" + strconv.Itoa(core.SrvPort))
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

		go handle_clt(client)
	}
}

func handshake_srv(client net.Conn) error {
	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	n, err := client.Read(buff)
	if err != nil {
		return err;
	}

	if string(buff[:n]) != core.Token {
		return errors.New("token err");
	}

	client.Write([]byte("ok"))

	return nil
}

func handle_clt(client net.Conn) {
	defer client.Close()

	err := handshake_srv(client)
	if err != nil {
		return
	}

	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	n, err := client.Read(buff)
	if err != nil {
		return;
	}

	fmt.Printf("new connect from %s to %s \n",client.RemoteAddr(), buff[:n])
	target, err := net.Dial("tcp", string(buff[:n]))
	if err != nil {
		return;
	}
	defer target.Close()

	client.Write([]byte("ok"))

	go func() {
		defer target.Close()
		defer client.Close()

		for {
			buff, n, err := core.ProxyRead(client)
			if err != nil {
				return
			}
			target.Write(buff[:n])
		}
	}()

	for{
		n,err := target.Read(buff)
		if err != nil {
			return
		}

		core.ProxyWrite(client, buff[:n])
	}
}