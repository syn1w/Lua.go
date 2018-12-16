package vm

import "testing"

func TestInt2Fb(t *testing.T) {
	for i := 0; i < 0xFF; i++ { // i is float point byte
		got := Int2fb(Fb2int(i))
		if got != i {
			t.Errorf("want: %d, got: %d", i, got)
		}
	}
}
