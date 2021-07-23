package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"goproxy/cipher"
	"goproxy/core"
	"goproxy/utils"
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

func (this *Socks5) Connected(c net.Conn) error {
	_, err := c.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	return err
}

func (this *Socks5) Process (c net.Conn) {
	defer c.Close()

	err := this.Auth(c)
	if err != nil {
		fmt.Printf("auth err: %s\n", err)
		return;
	}

	target, err := this.GetTarget(c)
	if err != nil {
		return
	}

	cip , err := cipher.GetDriver(utils.Ini)
	if err != nil {
		fmt.Printf("cipher err: %s\n", err)
		return
	}

	s, err := core.Connect(utils.Ini.GetString("srv", "host"), utils.Ini.GetString("srv", "port"), cip)
	if err != nil {
		fmt.Printf("Connect err: %s\n", err)
		return
	}
	defer s.Close()

	s.HandShakeClt()

	_, err = s.Send([]byte(target))
	if err != nil {
		return
	}

	buff := core.Pool.Get().([]byte)
	defer core.Pool.Put(buff)

	_, err = s.Recv(buff)
	if err != nil {
		return
	}

	fmt.Printf("new connect from %s to %s \n", c.RemoteAddr(), target)
	err = this.Connected(c)
	if err != nil {
		return
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