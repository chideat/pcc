package pig

import (
	"testing"
)

func TestUint64(t *testing.T) {
	_, err := Uint64(0)
	if err == nil {
		t.Fail()
	}

	id, err := Uint64(10)
	if id == 0 || err != nil {
		t.Fail()
	}
}

func BenchmarkUint64(b *testing.B) {
	for j := 0; j < b.N; j++ {
		Uint64(100)
	}
}
