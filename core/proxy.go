package core

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)


func Serve(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodConnect {
		handleHttps(w, r)
	} else {
		handleHttp(w, r)
	}
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	body := r.Body;
	str := make([]byte, 0);
	body.Read(str);
	var serverConn, dialErr = net.Dial("tcp", string(str));
	if dialErr != nil {
		return;
	}

	io.Copy(w, serverConn)
}



func handleHttps(w http.ResponseWriter, r *http.Request){
	destConn, err := net.DialTimeout("tcp", ProxyIP+":"+ProxyPort, 60*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	ProxySend(destConn, "c", []byte(r.RequestURI))
	msg, err := ProxyRecv(destConn)
	if err != nil {
		log.Println("proxy srv read err")
		return;
	}
	if msg.Op != "o" {
		log.Println("proxy srv handshak err")
		return;
	}
	log.Printf("proxy srv connect %s\n", msg.Data)

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	defer clientConn.Close()

	log.Println("proxy start")
	go ProxyTransfer(destConn, clientConn)
	for{
		msg, err := ProxyRecv(destConn)
		if err != nil {
			break;
		}
		switch msg.Op {
		case "p" :
			clientConn.Write([]byte(msg.Data))
		}
	}
	log.Println("proxy over")
}


func handleHttp(w http.ResponseWriter, r *http.Request){
	destConn, err := net.DialTimeout("tcp", ProxyIP+":"+ProxyPort, 60*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	host := r.Host
	if !strings.Contains(host, ":") {
		host += ":80"
	}

	log.Println(host)

	ProxySend(destConn, "c", []byte(host))
	msg, err := ProxyRecv(destConn)
	if err != nil {
		log.Println("proxy srv read err")
		return;
	}
	if msg.Op != "o" {
		log.Println("proxy srv handshak err")
		return;
	}
	log.Printf("proxy srv connect %s\n", msg.Data)

	var buff string
	buff += r.Method + " " + r.RequestURI + " " + "\r\n"

	for k,v := range(r.Header) {
		buff += k + ": "
		for id,n := range (v) {
			if id > 0 {
				buff += ", "
			}
			buff += n
		}
		buff += "\r\n"
	}

	buff += "\r\n"

	ProxySend(destConn, "p", []byte(buff))

	log.Printf("%s \n", buff)

	go func() {
		buff := make([]byte, 1024*10)
		n,err := r.Body.Read(buff)
		if err != nil {
			return
		}
		ProxySend(destConn, "p", buff[:n])
	}()

	for{
		msg, err := ProxyRecv(destConn)
		if err != nil {
			destConn.Close()
			break;
		}
		switch msg.Op {
		case "p" :
			w.Write([]byte(msg.Data))
		}
	}
	log.Println("proxy over")
}
