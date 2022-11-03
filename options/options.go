package options

import (
	"github.com/1uLang/libnet/encrypt"
	"time"
)

// Options Server初始化参数
type Options struct {
	EncryptMethod encrypt.MethodInterface // 数据加解密算法
	Timeout       time.Duration           // 连接读写超时时间
	PrivateKey    []byte                  // 加解密算法私钥
	PublicKey     []byte                  // 加解密算法公钥
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

// WithTimeout 设置TCP超时检查的间隔时间以及超时时间
func WithTimeout(timeout time.Duration) Option {
	return newFuncServerOption(func(o *Options) {
		if timeout < 0 {
			panic("timeoutTicker must greater than 0")
		}
		o.Timeout = timeout
	})
}

// WithPrivateKey 设置加解密私钥
func WithPrivateKey(privateKey []byte) Option {
	return newFuncServerOption(func(o *Options) {
		if len(privateKey) == 0 {
			panic("privateKey not be nil")
		}
		o.PrivateKey = privateKey
	})
}

// WithPublicKey 设置加解密公钥
func WithPublicKey(publicKey []byte) Option {
	return newFuncServerOption(func(o *Options) {
		if len(publicKey) < 0 {
			panic("publicKey not be nil")
		}
		o.PublicKey = publicKey
	})
}

func GetOptions(opts ...Option) *Options {
	options := &Options{}

	for _, o := range opts {
		o.apply(options)
	}
	return options
}
