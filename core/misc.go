package core

import (
	"bytes"
	"encoding/binary"
	"net"
	"strings"
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

func GetHost(buff []byte) string {
	headers := strings.Split(string(buff), "\r\n")
	for _, header := range headers {
		arr := strings.Split(header, ": ")
		if strings.ToLower(arr[0]) == "host" {
			if strings.Contains(arr[1], ":") {
				return arr[1]
			}
			return arr[1] + ":80"
		}
	}
	return ""
}

func GetMethod(buff []byte) string {
	headers := strings.Split(string(buff), "\r\n")
	arr := strings.Split(headers[0], " ")
	return arr[0]
}