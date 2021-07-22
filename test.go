package main

import (
	"fmt"
	"goproxy/core"
)

func main() {

	dst, _ := core.AesCtrEncrypt([]byte("zcs"), []byte("1234567887654321"))
	fmt.Printf("%s \n", dst)

	src, _ := core.AesCtrDecrypt(dst, []byte("1234567887654321"))
	fmt.Printf("%s \n", src)
}
