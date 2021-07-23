package core

import (
	"fmt"
	"goproxy/cipher"
	"net"
)

type Proxy struct {
	c net.Conn
	cipher cipher.Driver
}

func Listen(port string,  cipher cipher.Driver, handle func(Proxy)) error {
	s, err := net.Listen("tcp", ":" + port)
	if err != nil {
		return err
	}

	for {
		client, err := s.Accept()
		if err != nil {
			fmt.Printf("Accept failed: %v", err)
			break
		}

		p := Proxy{}
		p.c = client
		p.cipher = cipher

		go handle(p)
	}

	return err
}

func Connect(host string, port string, cipher cipher.Driver) (p Proxy, err error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return p, err
	}

	p.c = c
	p.cipher = cipher

	return p, nil
}

func (this *Proxy) SetCipher(driver cipher.Driver) {
	this.cipher = driver
}

func (this *Proxy) Send(b []byte) (int, error) {
	this.cipher.Encode(b, len(b))
	return this.c.Write(b)
}

func (this *Proxy) Recv(b []byte) (int, error) {
	n, err := this.c.Read(b)
	if err != nil {
		return n,err
	}

	this.cipher.Decode(b, n)
	return n,err
}

func (this *Proxy) Close() {
	this.c.Close()
}