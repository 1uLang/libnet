package main

import (
	"github.com/1uLang/libnet"
	"github.com/1uLang/libnet/encrypt"
	"github.com/1uLang/libnet/options"
	"time"
)

func main() {
	method, err := encrypt.NewMethodInstance("aes-256-cfb", encrypt.MagicKey, encrypt.MagicKey)
	if err != nil {
		panic(err)
	}
	svr := libnet.NewServe(":2439", new(Handle),
		options.WithEncryptMethod(method),
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
