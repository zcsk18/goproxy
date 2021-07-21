package core

import (
	"net"
	"strconv"
)

type Msg struct {
	Length int
	Op     string
	Data   []byte
}


func ProxySend(c net.Conn, op string, data []byte) error {
	buff := make([]byte, 512)

	length := len(data)
	copy(buff, IntToBytes(length))
	_, err := c.Write(buff[:strconv.IntSize]);
	if err != nil {
		return err
	}

	copy(buff, op);
	_, err = c.Write(buff[:1]);
	if err != nil {
		return err
	}

	Encode(&data, length)
	_, err = c.Write(data);
	if err != nil {
		return err
	}

	return nil
}


func ProxyRecv(c net.Conn) (*Msg, error){
	buff_len := make([]byte, strconv.IntSize)
	buff_op := make([]byte, 1)
	buff_data := make([]byte, 1024*10)
	msg := Msg{}
	length, err := c.Read(buff_len);
	if err != nil {
		return &msg, err;
	}
	if length <= 0 {
		return &msg, err;
	}
	msg.Length = BytesToInt(buff_len[:length])

	length, err = c.Read(buff_op);
	if err != nil {
		return &msg, err;
	}
	if length <= 0 {
		return &msg, err;
	}
	msg.Op = string(buff_op)

	length, err = c.Read(buff_data);
	if err != nil {
		return &msg, err;
	}
	if length <= 0 {
		return &msg, err;
	}
	Decode(&buff_data, length)
	msg.Data = buff_data[:length]
	return &msg, nil
}


func ProxyTransfer(dest net.Conn, src net.Conn) {
	buff_data := make([]byte, 1024*10)

	for {
		length, err := src.Read(buff_data)
		if err != nil {
			src.Close()
			dest.Close()
			return
		}
		if length <= 0 {
			src.Close()
			dest.Close()
			return
		}

		ProxySend(dest, "p", buff_data[:length])
	}
}