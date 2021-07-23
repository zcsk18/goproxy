package core

import (
	"errors"
	"fmt"
	"goproxy/cipher"
	"goproxy/utils"
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
		c, err := s.Accept()
		if err != nil {
			fmt.Printf("Accept failed: %v", err)
			break
		}

		p := Proxy{}
		p.c = c
		p.cipher = cipher

		go handle(p)
	}

	return err
}

func Connect(host string, port string, cipher cipher.Driver) (p Proxy, err error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
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

func (this *Proxy) RemoteAddr() net.Addr {
	return this.c.RemoteAddr()
}

func (this *Proxy) Close() {
	this.c.Close()
}

func (this *Proxy) HandShakeClt() error {
	_, err := this.Send([]byte(utils.GetIniParser().GetString("common", "token")))
	if err != nil {
		return err
	}

	buff := Pool.Get().([]byte)
	defer Pool.Put(buff)

	_, err = this.Recv(buff)
	if err != nil {
		return err
	}

	return nil
}

func (this *Proxy) HandShakeSrv() error {
	buff := Pool.Get().([]byte)
	defer Pool.Put(buff)

	n, err := this.Recv(buff)
	if err != nil {
		return err
	}

	if string(buff[:n]) != utils.GetIniParser().GetString("common", "token") {
		return errors.New("token err");
	}

	this.Send([]byte("ok"))
	return nil
}