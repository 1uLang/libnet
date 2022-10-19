package encrypt

const (
	MagicKey = "f1c8eafb543f03023e97b7be864a4e9b"
)

// 加密特殊信息
func MagicKeyEncode(data []byte) []byte {
	method, err := NewMethodInstance("aes-256-cfb", MagicKey, MagicKey[:16])
	if err != nil {
		return data
	}

	dst, err := method.Encrypt(data)
	if err != nil {
		return data
	}
	return dst
}

// 解密特殊信息
func MagicKeyDecode(data []byte) []byte {
	method, err := NewMethodInstance("aes-256-cfb", MagicKey, MagicKey[:16])
	if err != nil {
		return data
	}

	src, err := method.Decrypt(data)
	if err != nil {
		return data
	}
	return src
}
