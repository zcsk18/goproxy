package core

import (
	"net"
	"sync"
)

var Pool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024*100)
	},
}


func Encode(bufPtr *[]byte, num int) {
	buf := *bufPtr
	for i := 0; i < num; i++ {
		buf[i] = ^buf[i]
	}
}

func Decode(bufPtr *[]byte, num int) {
	buf := *bufPtr
	for i := 0; i < num; i++ {
		buf[i] = ^buf[i]
	}
}

func ProxyRead(c net.Conn, b []byte) (n int, err error) {
	n,err = c.Read(b)
	if err != nil {
		return n, err
	}

	Decode(&b, n)

	return n, err
}


func ProxyWrite(c net.Conn, b []byte) (n int, err error) {
	Encode(&b, n)

	n, err = c.Write(b)
	return n, err
}