// Package xchain @utils.go contains utilites for data structure manipulation.
// At the time of writing these are pertinent only to Parity data values, and
// might move to that package at some point.

package xchain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Uint64 uint64

func (u *Uint64) UnmarshalJSON(input []byte) error {
	if input[0] == '"' {
		uq, err := strconv.Unquote(strings.ToLower(string(input)))
		if err != nil {
			return err
		}
		input = []byte(uq)
	}
	ui, err := strconv.ParseUint(string(input), 0, 64)
	if err != nil {
		return err
	}
	*u = Uint64(ui)
	return nil

	return nil
}

func (u *Uint64) MarshalJSON() ([]byte, error) {
	x := hexutil.Uint64(uint64(*u))
	return json.Marshal(x)
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint64) MarshalText() ([]byte, error) {
	return hexutil.Uint64(u).MarshalText()
}

type BlockReward map[Uint64]*hexutil.Big

func (br *BlockReward) UnmarshalJSON(input []byte) error {
	if input[0] != '{' {
		var hb = new(hexutil.Big)
		if err := hb.UnmarshalJSON(input); err != nil {
			return err
		}
		zero := Uint64(0)
		*br = BlockReward{zero: hb}
		return nil
	}

	type BlockRewardMap map[string]string
	m := BlockRewardMap{}
	if err := json.Unmarshal(input, &m); err != nil {
		return err
	}
	var bbr = BlockReward{}
	for k, v := range m {
		var u Uint64
		err := u.UnmarshalJSON([]byte(k))
		if err != nil {
			return err
		}
		var hb = new(hexutil.Big)
		err = hb.UnmarshalJSON([]byte(strconv.Quote(v)))
		if err != nil {
			return err
		}
		bbr[u] = hb
	}
	*br = bbr
	return nil
}

type BTreeMap map[Uint64]*Uint64

func (btm *BTreeMap) UnmarshalJSON(input []byte) error {
	type IntermediateMap map[string]interface{}
	// m := make(map[string]interface{})
	m := IntermediateMap{}
	err := json.Unmarshal(input, &m)
	if err != nil {
		return err
	}
	var bbtm = BTreeMap{}
	for k, v := range m {
		var ku Uint64
		err := ku.UnmarshalJSON([]byte(k))
		if err != nil {
			return err
		}
		var vu Uint64
		vv, ok := v.(float64)
		if ok {
			vu = Uint64(vv)
		} else {
			vs, ok := v.(string)
			if ok {
				err = vu.UnmarshalJSON([]byte(vs))
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("could not assert btree map type")
			}
		}
		bbtm[ku] = &vu
	}
	*btm = bbtm
	return nil
}
