package protocol

import "net"

type Driver interface {
	Auth(c net.Conn) error
	GetTarget(c net.Conn) (string, error)
}