package options

import (
	"github.com/1uLang/libnet/encrypt"
	"time"
)

// Options Server初始化参数
type Options struct {
	EncryptMethod encrypt.MethodInterface
	Key           []byte
	Iv            []byte
	Timeout       time.Duration
}

type Option interface {
	apply(*Options)
}

type funcServerOption struct {
	f func(*Options)
}

func (fdo *funcServerOption) apply(do *Options) {
	fdo.f(do)
}

func newFuncServerOption(f func(*Options)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

// WithEncryptMethod 设置加解密方法
func WithEncryptMethod(encryptMethod encrypt.MethodInterface) Option {
	return newFuncServerOption(func(o *Options) {
		o.EncryptMethod = encryptMethod
	})
}

// WithEncryptMethodPublicKey 设置加解密方法公钥
func WithEncryptMethodPublicKey(key []byte) Option {
	return newFuncServerOption(func(o *Options) {
		if len(key) <= 0 {
			panic("encrypt public key must greater than 0")
		}
		o.Key = key
	})
}

// WithEncryptMethodPrivateKey 设置加解密方法私钥
func WithEncryptMethodPrivateKey(iv []byte) Option {
	return newFuncServerOption(func(o *Options) {
		i := len(iv)
		if i <= 0 {
			panic("encrypt private key must greater than 0")
		}
		o.Iv = iv
	})
}

// WithTimeout 设置TCP超时检查的间隔时间以及超时时间
func WithTimeout(timeout time.Duration) Option {
	return newFuncServerOption(func(o *Options) {
		if timeout <= 0 {
			panic("timeoutTicker must greater than 0")
		}

		o.Timeout = timeout
	})
}

func GetOptions(opts ...Option) *Options {
	options := &Options{}

	for _, o := range opts {
		o.apply(options)
	}
	if options.EncryptMethod != nil {
		if err := options.EncryptMethod.Init(options.Key, options.Iv); err != nil {
			panic(err)
		}
	}
	return options
}
func CheckOptions(opt *Options) error {

	if opt.EncryptMethod != nil {
		return opt.EncryptMethod.Init(opt.Key, opt.Iv)
	}
	return nil
}
