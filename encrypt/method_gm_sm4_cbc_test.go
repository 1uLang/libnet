package encrypt

import "testing"

func TestGMSM4CBCMethod_Encrypt(t *testing.T) {
	method, err := NewMethodInstance("gm-sm4-cbc", "abc", "123")
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("Hello, World")
	dst, err := method.Encrypt(src)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dst:", string(dst))

	src, err = method.Decrypt(dst)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("src:", string(src))
}

func TestGMSM4CBCMethod_Encrypt2(t *testing.T) {
	method, err := NewMethodInstance("gm-sm4-cbc", "abc", "123")
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("Hello, World")
	dst, err := method.Encrypt(src)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dst:", string(dst))

	src, err = method.Decrypt(dst)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("src:", string(src))
}
