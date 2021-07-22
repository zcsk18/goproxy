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


func ProxyRead(c net.Conn, b []byte) (n int, err error) {
	n,err = c.Read(b)
	if err != nil {
		return n, err
	}

	return n, err
}


func ProxyWrite(c net.Conn, b []byte) (n int, err error) {
	n, err = c.Write(b)
	return n, err
}