package core

import (
	"bytes"
	"encoding/binary"
	"net"
)

func CheckAdress(adress string) bool{
	_, err := net.ResolveTCPAddr("tcp", adress)
	if err != nil{
		return false
	}
	return true
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


func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}