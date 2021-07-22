package core

import (
	"net"
	"strconv"
	"sync"
)

type Msg struct {
	Action string
	Data []byte
}

var Pool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024*100)
	},
}

func ProxySend(c net.Conn, action string, data []byte) error {
	buff := make([]byte, 512)

	length := len(data)
	copy(buff, IntToBytes(length))
	_, err := c.Write(buff[:strconv.IntSize]);
	if err != nil {
		return err
	}

	copy(buff, action);
	_, err = c.Write(buff[:1]);
	if err != nil {
		return err
	}

	//Encode(&data, len(data))

	totle := 0
	for {
		n, err := c.Write(data[totle:]);
		if err != nil {
			return err
		}
		totle += n
		if (totle == len(data)) {
			break
		}
	}

	return nil
}


func ProxyRecv(c net.Conn) (*Msg, error){
	buff_len := make([]byte, strconv.IntSize)
	buff_op := make([]byte, 1)
	buff_data := make([]byte, 1024*100)
	msg := Msg{}
	length, err := c.Read(buff_len);
	if err != nil {
		return &msg, err;
	}
	if length <= 0 {
		return &msg, err;
	}

	length, err = c.Read(buff_op);
	if err != nil {
		return &msg, err;
	}
	if length <= 0 {
		return &msg, err;
	}
	msg.Action = string(buff_op)

	length, err = c.Read(buff_data);
	if err != nil {
		return &msg, err;
	}

	//Decode(&buff_data, length)
	msg.Data = buff_data[:length]

	return &msg, nil
}