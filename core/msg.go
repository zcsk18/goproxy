package core

import (
	"net"
	"sync"
)

var Pool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024*500)
	},
}

const MaxValue  = 255
const Dis = 4


func Encode(b []byte, n int, dis byte) ([]byte) {
	for i:=0; i<n ;i++  {
		if (b[i] + dis > MaxValue) {
			b[i] = MaxValue - b[i] + (dis - 1)
		} else {
			b[i] += dis
		}
	}
	return b;
}


func Decode(b []byte, n int, dis byte) ([]byte) {
	for i:=0; i<n ;i++  {
		if (b[i] - dis < 0) {
			b[i] = MaxValue - (dis-1) - b[i]
		} else {
			b[i] -= dis
		}
	}
	return b;
}



func ProxyRead(c net.Conn) ([]byte, int, error) {
	buff := Pool.Get().([]byte)
	defer Pool.Put(buff)

	n,err := c.Read(buff)
	if err != nil {
		return buff, n, err
	}

	buff = Decode(buff, n, Dis)
	return buff, n, err
}


func ProxyWrite(c net.Conn, b []byte) (int, error) {
	b = Encode(b, len(b), Dis)
	n, err := c.Write(b)
	return n, err
}