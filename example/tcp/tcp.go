package main

import (
	"github.com/1uLang/libnet"
	"github.com/1uLang/libnet/encrypt"
	"github.com/1uLang/libnet/example"
	"github.com/1uLang/libnet/options"
	"time"
)

func main() {
	svr, err := libnet.NewServe(":2439", new(example.Handle),
		options.WithEncryptMethod(new(encrypt.AES256CFBMethod)),
		options.WithEncryptMethodPublicKey([]byte(encrypt.MagicKey)),
		options.WithEncryptMethodPrivateKey([]byte(encrypt.MagicKey[:16])),
		options.WithTimeout(5*time.Second),
	)
	if err != nil {
		panic(err)
	}
	err = svr.RunTCP()
	if err != nil {
		panic(err)
	}
}
