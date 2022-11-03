package workers

import "testing"

func TestGet(t *testing.T) {
	for i := 0; i < 40; i++ {
		t.Log(Get())
	}
}
