package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	encrypt2 "github.com/1uLang/libnet/encrypt"
	message2 "github.com/1uLang/libnet/example/message"
	"github.com/1uLang/libnet/message"
	"io/ioutil"
	"log"
	"net"
	"syscall"
	"time"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	port        = flag.String("port", "2439", "server port")
	protocol    = flag.String("proto", "tcp", "server type tcp / udp")
	connections = flag.Int("conn", 1, "number of tcp connections")
	caCrt       = flag.String("ca", "/data/libnet/certs/ca.crt", "tls ca cert")
	crt         = flag.String("crt", "/data/libnet/certs/client.crt", "tls client cert file")
	key         = flag.String("key", "/data/libnet/certs/client.key", "tls client key file")
	encrypt     = flag.String("encrypt", "", "set send or recv buffer encrypt method")
)

func main() {
	flag.Parse()

	setLimit()

	addr := *ip + ":" + *port
	var enc encrypt2.MethodInterface
	var err error
	var conns []net.Conn
	var buffer = message.NewBuffer(message2.CheckHeader)
	var msg = message2.Message{}
	buffer.OnMessage(func(msg message.MessageI) {
		fmt.Println("recv msg : ", string(msg.GetData()))
	})
	if *encrypt != "" {
		enc, err = encrypt2.NewMethodInstance(*encrypt, encrypt2.MagicKey, encrypt2.MagicKey[:16])
		if err != nil {
			log.Fatalf("set encrypt method errrr: %v", err)
		}
	}
	log.Printf("连接到 %s", addr)

	for i := 0; i < *connections; i++ {
		var c net.Conn
		if *protocol == "tls" {
			cert, err := ioutil.ReadFile(*caCrt)
			if err != nil {
				log.Fatalf("could not open certificate file: %v", err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(cert)

			certificate, err := tls.LoadX509KeyPair(*crt, *key)
			if err != nil {
				log.Fatalf("could not load certificate: %v", err)
			}

			// Create a HTTPS client and supply the created CA pool and certificate
			tlsConfig := &tls.Config{
				RootCAs:            caCertPool,
				ClientCAs:          caCertPool,
				Certificates:       []tls.Certificate{certificate},
				ClientAuth:         tls.RequireAndVerifyClientCert,
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true,
			}

			c, err = tls.Dial("tcp", addr, tlsConfig)
		} else {
			c, err = net.DialTimeout(*protocol, addr, 10*time.Second)
		}

		if err != nil {
			fmt.Println("failed to connect", i, err)
			i--
			continue
		}
		conns = append(conns, c)
		time.Sleep(time.Millisecond)
	}

	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()

	log.Printf("完成初始化 %d 连接", len(conns))

	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 5
	}

	for {
		buf := make([]byte, 1024)
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			log.Printf("连接 %d 发送数据", i)

			msg.Data = []byte("hello world\r\n")
			if enc != nil {
				enSend, err := enc.Encrypt(msg.Marshal())
				if err != nil {
					log.Fatalf("encrypt encode error %s", err)
				}
				n, err := conn.Write(enSend)
				log.Println(" send to server : hello world ", n, err)
			} else {
				n, err := conn.Write(msg.Data)
				log.Println(" send to server : hello world ", n, err)
			}
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Println(" recv from server err  ", err)
			}
			if *protocol != "udp" {
				var recv []byte

				if enc != nil {
					recv, err = enc.Decrypt(buf[:n])
					if err != nil {
						log.Println(" recv from server byte decode error  ", err)
						continue
					}
				} else {
					recv = buf[:n]
				}
				buffer.Write(recv)
			}
		}
	}
	//select{}
}
func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
}
