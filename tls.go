package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var caCertFile = "./cacert.pem"
var svrCertFile = "./svrcert.pem"
var svrKeyFile = "./svrkey.pem"

func main() {
	var addr string
	var listenMode bool
	flag.StringVar(&addr, "addr", "localhost:20012", "Address")
	flag.BoolVar(&listenMode, "l", false, "Listen mode")
	flag.StringVar(&caCertFile, "ca", "/etc/pki/CA/cacert.pem", "CA certification")
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
			reader := bufio.NewReader(conn)
			writer := bufio.NewWriter(conn)
			go func() {
				log.Printf("Accept connection %v\n", conn.RemoteAddr())
				defer func() {
					log.Printf("Closing %v\n", conn.RemoteAddr())
					conn.Close()
				}()

				log.Println("Begin to process connection")
				for {
					var rbytes []byte
					var err error
					rbytes, err = reader.ReadBytes('\n')
					if err != nil {
						log.Printf("Read err: %v\n", err)
						conn.Close()
						break
					}
					log.Printf("Read request: %v", string(rbytes))

					var ss = fmt.Sprintf("World %v\n", time.Now().String())
					_, err = writer.Write([]byte(ss))
					if err != nil {
						log.Printf("WriteString fail: %v\n", err)
						break
					}
					err = writer.Flush()
					if err != nil {
						log.Printf("WriteString fail: %v\n", err)
						break
					}
					log.Printf("Send response: %v", ss)
				}
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
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)
		for {
			var err error
			var rbytes []byte
			var ss = fmt.Sprintf("Hello %v\n", time.Now().String())
			_, err = writer.Write([]byte(ss))
			if err != nil {
				log.Printf("Write err: %v\n", err)
				conn.Close()
				break
			}
			err = writer.Flush()
			if err != nil {
				log.Printf("Flush err: %v\n", err)
				conn.Close()
				break
			}
			log.Printf("Send request: %v", ss)

			rbytes, err = reader.ReadBytes('\n')
			if err != nil {
				log.Printf("Read err: %v\n", err)
				conn.Close()
				break
			}
			log.Printf("Read response: %v", string(rbytes))

			time.Sleep(1 * time.Second)
		}
	}
}
