package main

import (
	"flag"
	"goproxy/core"
	"goproxy/http"
	"log"
	"os"
)


func init() {
	f, _ := os.OpenFile("goproxy_clt" + ".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	log.SetOutput(f)
	log.SetOutput(os.Stdout)
}


func main() {
	var listenAdress string
	port := core.LocalPort

	flag.StringVar(&listenAdress, "L", "0.0.0.0:"+port, "listen address.eg: 127.0.0.1:"+port)
	flag.Parse()

	if !core.CheckAdress(listenAdress){
		log.Fatal("-L listen address format incorrect.Please check it")
	}

	http.Serve(listenAdress)
}