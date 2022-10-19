package encrypt

const (
	encryptMethodRaw = iota
	encryptMethodAES128CFB
	encryptMethodAES192CFB
	encryptMethodAES256CFB
	encryptMethodGMSM2ECC
	encryptMethodGMSM3SUM
	encryptMethodGMSM4CBC
)

type MethodInterface interface {
	// 初始化
	Init(key []byte, iv []byte) error

	// 加密
	Encrypt(src []byte) (dst []byte, err error)

	// 解密
	Decrypt(dst []byte) (src []byte, err error)

	// 加密方式ID
	Method() uint8
}
