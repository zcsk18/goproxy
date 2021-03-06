package main

import (
	"fmt"
	"goproxy/cipher"
	"goproxy/core"
	"goproxy/utils"
	"net"
)

func main() {
	cip , err := cipher.GetDriver(utils.GetIniParser())
	if err != nil {
		return
	}

	core.Listen(utils.GetIniParser().GetString("srv", "port"), cip, process)
}

func process(c core.Proxy) {
	defer c.Close()

	err := c.HandShakeSrv()
	if err != nil {
		return
	}

	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	n, err := c.Recv(buff)
	if err != nil {
		return;
	}
	fmt.Printf("new connect from %s to %s \n",c.RemoteAddr(), buff[:n])

	target, err := net.Dial("tcp", string(buff[:n]))
	if err != nil {
		return;
	}
	defer target.Close()

	c.Send([]byte("ok"))

	go func() {
		defer target.Close()
		defer c.Close()

		buff := core.Pool.Get().([]byte)
		defer core.Pool.Put(buff)

		for {
			n, err := c.Recv(buff)
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

		c.Send(buff[:n])
	}
}