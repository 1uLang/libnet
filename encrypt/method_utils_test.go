package encrypt

import (
	"crypto/rand"
	"fmt"
	"github.com/ZZMarquis/gm/sm2"
	"testing"
)

func TestFindMethodInstance(t *testing.T) {
	t.Log(NewMethodInstance("a", "b", ""))
	t.Log(NewMethodInstance("aes-256-cfb", "123456", ""))
}

func TestSm2(t *testing.T) {

	pri, pub, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(pri.GetRawBytes()))
	fmt.Println(string(pub.GetRawBytes()))
}
