package encrypt

import (
	"errors"
	"reflect"
)

var (
	encryptKey = ""
	encryptIv  = ""
)

var methods = map[string]reflect.Type{
	"raw":         reflect.TypeOf(new(RawMethod)).Elem(),
	"aes-128-cfb": reflect.TypeOf(new(AES128CFBMethod)).Elem(),
	"aes-192-cfb": reflect.TypeOf(new(AES192CFBMethod)).Elem(),
	"aes-256-cfb": reflect.TypeOf(new(AES256CFBMethod)).Elem(),
	"gm-sm2-ecc":  reflect.TypeOf(new(GMSM2ECCMethod)).Elem(),
	"gm-sm3-sum":  reflect.TypeOf(new(GMSM3SUMMethod)).Elem(),
	"gm-sm4-cbc":  reflect.TypeOf(new(GMSM4CBCMethod)).Elem(),
}

func Init(key string, iv string) {
	encryptKey, encryptIv = key, iv
}
func NewMethod(method string) (MethodInterface, error) {
	valueType, ok := methods[method]
	if !ok {
		return nil, errors.New("method '" + method + "' not found")
	}
	instance, ok := reflect.New(valueType).Interface().(MethodInterface)
	if !ok {
		return nil, errors.New("method '" + method + "' must implement MethodInterface")
	}
	return instance, nil
}
func NewMethodInstance(method string, key string, iv string) (MethodInterface, error) {
	instance, err := NewMethod(method)
	if err != nil {
		return nil, err
	}
	err = instance.Init([]byte(key), []byte(iv))
	return instance, err
}
func GetMethodInstance(id uint8) (MethodInterface, error) {
	method := ""
	switch id {
	case encryptMethodRaw:
		method = "raw"
	case encryptMethodAES128CFB:
		method = "aes-128-cfb"
	case encryptMethodAES192CFB:
		method = "aes-192-cfb"
	case encryptMethodAES256CFB:
		method = "aes-256-cfb"
	case encryptMethodGMSM2ECC:
		method = "gm-sm2-ecc"
	case encryptMethodGMSM3SUM:
		method = "gm-sm3-sum"
	case encryptMethodGMSM4CBC:
		method = "gm-sm4-cbc"
	}
	return NewMethodInstance(method, encryptKey, encryptIv)
}
func RecoverMethodPanic(err interface{}) error {
	if err != nil {
		s, ok := err.(string)
		if ok {
			return errors.New(s)
		}

		e, ok := err.(error)
		if ok {
			return e
		}

		return errors.New("unknown error")
	}
	return nil
}
