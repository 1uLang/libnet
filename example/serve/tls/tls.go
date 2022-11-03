package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/1uLang/libnet"
	"github.com/1uLang/libnet/options"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	svr := libnet.NewServe(":2439", new(Handle),
		//options.WithEncryptMethod(new(encrypt.AES256CFBMethod)),
		//options.WithEncryptMethodPublicKey([]byte(encrypt.MagicKey)),
		//options.WithEncryptMethodPrivateKey([]byte(encrypt.MagicKey[:16])),
		options.WithTimeout(5*time.Second),
	)
	caCertFile, err := ioutil.ReadFile("/data/libnet/certs/ca.crt")
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)

	certificate, err := tls.LoadX509KeyPair("/data/libnet/certs/server.crt", "/data/libnet/certs/server.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}
	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{certificate},
		ClientCAs:                caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		InsecureSkipVerify: true,
	}
	err = svr.RunTLS(tlsConfig)
	if err != nil {
		panic(err)
	}
}
