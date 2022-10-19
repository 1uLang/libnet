package encrypt

import (
	"github.com/ZZMarquis/gm/sm2"
)

type GMSM2ECCMethod struct {
	iv  []byte
	key []byte
}

func (this *GMSM2ECCMethod) Init(pri, pub []byte) error {
	_, err := sm2.RawBytesToPublicKey(pub)
	if err != nil {
		return err
	}
	_, err = sm2.RawBytesToPrivateKey(pri)
	if err != nil {
		return err
	}
	this.key = pri
	this.iv = pub
	return nil
}

func (this *GMSM2ECCMethod) Encrypt(in []byte) (dst []byte, err error) {

	pub, err := sm2.RawBytesToPublicKey(this.iv)
	if err != nil {
		return
	}
	cipherText, err := sm2.Encrypt(pub, in, sm2.C1C2C3)
	if err != nil {
		return
	}

	return cipherText, nil
}

func (this *GMSM2ECCMethod) Decrypt(out []byte) (in []byte, err error) {
	pri, err := sm2.RawBytesToPrivateKey(this.key)
	if err != nil {
		return
	}

	plainText, err := sm2.Decrypt(pri, out, sm2.C1C2C3)
	if err != nil {
		return
	}
	return plainText, nil
}
func (this *GMSM2ECCMethod) Method() uint8 {
	return encryptMethodGMSM2ECC
}
