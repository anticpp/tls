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

var caCertPem = "./cacert.pem"
var svrCertPem = "./svrcert.pem"
var svrKeyPem = "./svrkey.pem"

func main() {
	var addr string
	var listen bool
	flag.StringVar(&addr, "addr", "localhost:20012", "Address")
	flag.BoolVar(&listen, "l", false, "Listen mode")
	flag.StringVar(&caCertPem, "ca", "/etc/pki/CA/cacert.pem", "CA certification")
	flag.Parse()

	if listen {
		cert, err := tls.LoadX509KeyPair(svrCertPem, svrKeyPem)
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
				defer func() {
					log.Printf("Closing %v\n", conn.RemoteAddr())
					conn.Close()
				}()

				log.Printf("Accept connection %v\n", conn.RemoteAddr())
				reader := bufio.NewReader(conn)
				writer := bufio.NewWriter(conn)

				var rbuf []byte
				var err error
				rbuf, err = reader.ReadBytes('\n')
				if err != nil {
					log.Printf("Read err: %v\n", err)
					conn.Close()
					return
				}
				log.Printf("Read request: %v", string(rbuf))

				var sbuf = fmt.Sprintf("World %v\n", time.Now().String())
				_, err = writer.Write([]byte(sbuf))
				if err != nil {
					log.Printf("Write err: %v\n", err)
					return
				}
				writer.Flush()
				log.Printf("Send response: %v", sbuf)
			}()
		}
	} else {
		fmt.Printf("Reading CA certificate %v\n", caCertPem)
		cert, err := ioutil.ReadFile(caCertPem)
		if err != nil {
			log.Fatalf("Read CaCertPem fail: %v\n", err)
		}

		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(cert)
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
		for {
			defer func() {
				log.Printf("Closing %v\n", conn.RemoteAddr())
				conn.Close()
			}()

			reader := bufio.NewReader(conn)
			writer := bufio.NewWriter(conn)

			var err error
			var sbuf = fmt.Sprintf("Hello %v\n", time.Now().String())
			_, err = writer.Write([]byte(sbuf))
			if err != nil {
				log.Printf("Write err: %v\n", err)
				break
			}
			writer.Flush()
			log.Printf("Send request: %v", sbuf)

			var rbuf []byte
			rbuf, err = reader.ReadBytes('\n')
			if err != nil {
				log.Printf("Read err: %v\n", err)
				conn.Close()
				break
			}
			log.Printf("Read response: %v", string(rbuf))
			break
		}
	}
}
