package xchain

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestBalanceBigSetString(t *testing.T) {
	v := "1940000000000000000000"
	b, ok := ParseBig256(v)
	if !ok {
		t.Fail()
	}
	t.Log(b)
}

func TestBlockRewardMarshaling(t *testing.T) {
	input := []byte(`{
			  "0x0": "0x4563918244F40000",
			  "0x5": "0x29A2241AF62C0000"
		  }`)

	br := BlockReward{}
	err := json.Unmarshal(input, &br)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(spew.Sdump(br))
}
