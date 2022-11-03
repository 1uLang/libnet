package options

import (
	"github.com/1uLang/libnet/encrypt"
	"time"
)

// Options Server初始化参数
type Options struct {
	EncryptMethod encrypt.MethodInterface // 数据加解密算法
	Timeout       time.Duration           // 连接读写超时时间
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

func GetOptions(opts ...Option) *Options {
	options := &Options{}

	for _, o := range opts {
		o.apply(options)
	}
	return options
}
