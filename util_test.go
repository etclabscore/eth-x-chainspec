package xchain

import "testing"

func TestBalanceBigSetString(t *testing.T) {
	v := "1940000000000000000000"
	b, ok := ParseBig256(v)
	if !ok {
		t.Fail()
	}
	t.Log(b)
}
