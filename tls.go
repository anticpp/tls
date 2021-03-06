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

var caCertPem = "/etc/pki/CA/cacert.pem"
var svrCertPem = "./svrcert.pem"
var svrKeyPem = "./svrkey.pem"
var clientCertPem = "./cltcert.pem"
var clientKeyPem = "./cltkey.pem"

func main() {
	var addr string
	var listen bool
	var mtls bool
	var serverName string
	flag.StringVar(&addr, "addr", "localhost:20012", "Address.")
	flag.BoolVar(&listen, "l", false, "Listen mode.")
	flag.StringVar(&serverName, "sn", "www.example.com", "Server name.")
	flag.BoolVar(&mtls, "mtls", false,
		"mTLS on. Server must verify client certificate. Client side should provide client certificate.")
	flag.StringVar(&caCertPem, "ca", "/etc/pki/CA/cacert.pem", "CA certification")
	flag.Parse()

	log.Printf("Using CA certificate %v\n", caCertPem)
	cert, err := ioutil.ReadFile(caCertPem)
	if err != nil {
		log.Fatalf("Read CaCertPem fail: %v\n", err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		log.Fatalf("AppendCertsFromPEM fail: %v\n", err)
	}
	if listen {
		log.Printf("Using server key %v\n", svrKeyPem)
		log.Printf("Using server certificate %v\n", svrCertPem)
		cert, err := tls.LoadX509KeyPair(svrCertPem, svrKeyPem)
		if err != nil {
			log.Fatalf("LoadX509KeyPair fail; %v\n", err)
		}
		c := &tls.Config{
			// Server certificate
			Certificates: []tls.Certificate{cert},
		}
		if mtls {
			// Verify client certificate
			// Using ClientCAs
			log.Println("mTLS on. Client certificate must be verified.")
			c.ClientCAs = roots
			c.ClientAuth = tls.RequireAndVerifyClientCert
		}
		log.Printf("Listening %v\n", addr)
		l, err := tls.Listen("tcp", addr, c)
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
				log.Printf("%v", string(rbuf))

				var sbuf = fmt.Sprintf("[%v] Server World\n",
					time.Now().Format(time.UnixDate))
				_, err = writer.Write([]byte(sbuf))
				if err != nil {
					log.Printf("Write err: %v\n", err)
					return
				}
				writer.Flush()
				log.Printf("%v", sbuf)
			}()
		}
	} else {
		c := &tls.Config{
			RootCAs:    roots,
			ServerName: serverName,
		}
		if mtls {
			// Client certificate
			log.Println("mTLS on. Set client certificate.")
			log.Printf("Using client key %v\n", clientKeyPem)
			log.Printf("Using client certificate %v\n", clientCertPem)
			cert, err := tls.LoadX509KeyPair(clientCertPem, clientKeyPem)
			if err != nil {
				log.Fatalf("LoadX509KeyPair fail; %v\n", err)
			}

			c.Certificates = []tls.Certificate{cert}
		}
		log.Printf("Dialing %v\n", addr)
		log.Printf("Server name %v\n", serverName)
		conn, err := tls.Dial("tcp", addr, c)
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
			var sbuf = fmt.Sprintf("[%v] Client Hello\n",
				time.Now().Format(time.UnixDate))
			_, err = writer.Write([]byte(sbuf))
			if err != nil {
				log.Printf("Write err: %v\n", err)
				break
			}
			writer.Flush()
			log.Printf("%v", sbuf)

			var rbuf []byte
			rbuf, err = reader.ReadBytes('\n')
			if err != nil {
				log.Printf("Read err: %v\n", err)
				conn.Close()
				break
			}
			log.Printf("%v", string(rbuf))
			break
		}
	}
}
