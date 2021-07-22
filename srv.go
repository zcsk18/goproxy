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
	msg, err := core.ProxyRecv(client)
	if err != nil {
		client.Close()
		return err;
	}
	if string(msg.Data) != core.Token {
		client.Close()
		return errors.New("token err");
	}
	core.ProxySend(client, "T", []byte("ok"))

	return nil
}

func handle_clt(client net.Conn) {
	defer client.Close()

	handshake_srv(client)

	msg, err := core.ProxyRecv(client)
	if err != nil {
		return;
	}
	fmt.Printf("new connect from %s to %s \n",client.RemoteAddr(), msg.Data)
	target, err := net.Dial("tcp", string(msg.Data))
	defer target.Close()
	if err != nil {
		return;
	}
	core.ProxySend(client, "C", []byte("ok"))

	go func() {
		defer target.Close()
		defer client.Close()

		buff := core.Pool.Get().([]byte)
		defer core.Pool.Put(buff)

		for {
			n, err := client.Read(buff)
			if err != nil {
				return
			}
			core.Decode(&buff, n)
			target.Write(buff[:n])
		}
	}()

	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	for{
		n,err := target.Read(buff)
		if err != nil {
			return
		}
		core.Encode(&buff, n)
		client.Write(buff[:n])
	}
}