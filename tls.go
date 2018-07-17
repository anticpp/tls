package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
)

const caCertFile = "./cacert.pem"
const svrCertFile = "./svrcert.pem"
const svrKeyFile = "./svrkey.pem"

func main() {
	var addr string
	var listenMode bool
	flag.StringVar(&addr, "addr", "localhost:20012", "Address")
	flag.BoolVar(&listenMode, "l", false, "Listen mode")
	flag.Parse()

	if listenMode {
		cert, err := tls.LoadX509KeyPair(svrCertFile, svrKeyFile)
		if err != nil {
			log.Fatalf("LoadX509KeyPair fail; %v\n", err)
		}
		log.Printf("Listening %v\n", addr)
		l, err := tls.Listen("tcp", addr, &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			log.Fatalf("Listen fail: %v\n", err)
		}

		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("Accept fail: %v", err)
				continue
			}
			go func() {
				log.Printf("Accept connection %v\n", conn.RemoteAddr())
				defer func() {
					log.Printf("Closing %v\n", conn.RemoteAddr())
					conn.Close()
				}()
				n, err := conn.Write([]byte("Hello World"))
				if err != nil {
					log.Printf("WriteString fail: %v\n", err)
					return
				}
				log.Printf("WriteString success bytes %v\n", n)
			}()
		}
	} else {
		caCertPEM, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			log.Fatalf("Read CaCertPem fail: %v\n", err)
		}

		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(caCertPEM)
		if !ok {
			log.Fatalf("AppendCertsFromPEM fail: %v\n", err)
		}

		log.Printf("Dialing %v\n", addr)
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			RootCAs: roots,
		})
		if err != nil {
			log.Fatalf("Dial fail addr %v, %v", addr, err)
		}

		log.Printf("Connect success\n")
		tmpBuf := make([]byte, 1024)
		for {
			n, err := conn.Read(tmpBuf)
			if err != nil {
				log.Printf("Read err: %v\n", err)
				conn.Close()
				break
			}
			s := string(tmpBuf[0:n])
			log.Printf("ReadString: %v\n", s)
		}
	}
}
