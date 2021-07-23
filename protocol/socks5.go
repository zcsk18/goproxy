package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type Socks5 struct {
}


func (this *Socks5) Auth(c net.Conn) error {
	buf := make([]byte, 256)

	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(c, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}

	// 读取 METHODS 列表
	n, err = io.ReadFull(c, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}

	//无需认证
	n, err = c.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write rsp: " + err.Error())
	}

	return nil
}

func (this *Socks5) GetTarget(c net.Conn) (string, error){
	buf := make([]byte, 256)

	n, err := io.ReadFull(c, buf[:4])
	if n != 4 {
		return "", errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return "", errors.New("invalid ver/cmd")
	}

	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(c, buf[:4])
		if n != 4 {
			return "", errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case 3:
		n, err = io.ReadFull(c, buf[:1])
		if n != 1 {
			return "", errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(c, buf[:addrLen])
		if n != addrLen {
			return "", errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])

	case 4:
		return "", errors.New("IPv6: no supported yet")

	default:
		return "", errors.New("invalid atyp")
	}

	n, err = io.ReadFull(c, buf[:2])
	if n != 2 {
		return "", errors.New("read port: " + err.Error())
	}

	port := binary.BigEndian.Uint16(buf[:2])
	destAddrPort := fmt.Sprintf("%s:%d", addr, port)

	return destAddrPort, nil
}