package proto

import "net"

type Driver interface {
	Process (c net.Conn)
}