package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

const caCrt = "/Users/1usir/works/zero-trust/libnet/certs/CertAuth.crt"
const crt = "/Users/1usir/works/zero-trust/libnet/certs/client.crt"
const key = "/Users/1usir/works/zero-trust/libnet/certs/client.key"

func main() {
	cert, err := ioutil.ReadFile(caCrt)
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	// Create a HTTPS client and supply the created CA pool and certificate
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	c, err := tls.Dial("tcp", "192.168.1.118:2439", tlsConfig)
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("send hello")
		_, err := c.Write([]byte("hello"))
		if err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
	c.Close()
	return
}
