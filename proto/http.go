package proto

import (
	"fmt"
	"goproxy/cipher"
	"goproxy/core"
	"goproxy/utils"
	"net"
	"strings"
)

type Http struct {
	tp string
	target string
}

func (this *Http) Connected(c net.Conn) error {
	if this.tp == "https" {
		_, err := c.Write([]byte("HTTP/1.0 200 Connetion established\r\n\r\n"))
		return err
	}
	return nil
}


func (this *Http) Process (c net.Conn) {
	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	n, err := c.Read(buff)
	if err != nil {
		return
	}

	headers := strings.Split(string(buff[:n]), "\r\n")
	header := strings.Split(headers[0], " ")

	if header[0] == "CONNECT" {
		this.tp = "https"
		this.target = header[1];
	} else {
		this.tp = "http"
		for _, header := range headers {
			line := strings.Split(header, ": ")
			if line[0] == "Host" {
				this.target = line[1]
			}
		}
	}

	if !strings.Contains(this.target, ":") {
		this.target += ":80"
	}

	cip , err := cipher.GetDriver(utils.GetIniParser())
	if err != nil {
		fmt.Printf("cipher err: %s\n", err)
		return
	}

	s, err := core.Connect(utils.GetIniParser().GetString("srv", "host"), utils.GetIniParser().GetString("srv", "port"), cip)
	if err != nil {
		fmt.Printf("Connect err: %s\n", err)
		return
	}
	defer s.Close()

	s.HandShakeClt()

	_, err = s.Send([]byte(this.target))
	if err != nil {
		return
	}

	_, err = s.Recv(buff)
	if err != nil {
		return
	}

	fmt.Printf("new connect from %s to %s \n", c.RemoteAddr(), this.target)
	err = this.Connected(c)
	if err != nil {
		return
	}
	if this.tp == "http" {
		s.Send(buff[:n])
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











