package main

import (
	"github.com/1uLang/libnet"
	"github.com/1uLang/libnet/options"
	"time"
)

func main() {
	svr := libnet.NewServe(":2439", new(Handle),
		//options.WithEncryptMethod(new(encrypt.AES256CFBMethod)),
		//options.WithEncryptMethodPublicKey([]byte(encrypt.MagicKey)),
		//options.WithEncryptMethodPrivateKey([]byte(encrypt.MagicKey[:16])),
		options.WithTimeout(5*time.Second),
	)
	err := svr.RunUDP()
	if err != nil {
		panic(err)
	}
	select {}
}
