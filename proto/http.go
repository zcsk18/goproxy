package proto

import (
	"goproxy/core"
	"io"
	"net"
	"strings"
)

type Http struct {
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
		this.ProcessHttps(c, buff, n, header[1])
	} else {
		for _, header := range headers {
			line := strings.Split(header, ": ")
			if line[0] == "Host" {
				this.ProcessHttp(c, buff, n, line[1])
			}
		}
	}
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func (this *Http) ProcessHttps(c net.Conn, buff []byte, n int, target string) {
	s, err := net.Dial("tcp", target)
	if err != nil {
		return
	}

	c.Write([]byte("HTTP/1.0 200 Connetion established\r\n\r\n"))

	go transfer(s, c)
	go transfer(c, s)
}

func (this *Http) ProcessHttp (c net.Conn, buff []byte, n int, target string) {
	if !strings.Contains(target, ":") {
		target += ":80"
	}

	s, err := net.Dial("tcp", target)
	if err != nil {
		return
	}

	s.Write(buff[:n])

	go transfer(s, c)
	go transfer(c, s)
}












