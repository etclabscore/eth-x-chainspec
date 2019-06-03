// Package tests gets its own package because go-ethereums/tests
// and ./parity created an unallowed looping import cycle.

package tests

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/etclabscore/eth-x-chainspec/parity"
	"github.com/ethereum/go-ethereum/tests"
)

func TestXTestChainspecs(t *testing.T) {
	forks := tests.Forks
	testsChainspecsMap := map[string]string{
		"Frontier":          "frontier_test.json",
		"Homestead":         "homestead_test.json",
		"EIP150":            "eip150_test.json",
		"EIP158":            "eip161_test.json",
		"Byzantium":         "byzantium_test.json",
		"Constantinople":    "constantinople_test.json",
		"ConstantinopleFix": "st_peters_test.json",
	}
	for f, cf := range testsChainspecsMap {
		b, err := ioutil.ReadFile(filepath.Join("../parity/chainspecs", cf))
		if err != nil {
			t.Fatal(err)
		}

		pc := &parity.Config{}
		err = json.Unmarshal(b, &pc)
		if err != nil {
			t.Fatal(err)
		}

		mg := pc.ToMultiGethGenesis()

		if !reflect.DeepEqual(forks[f], mg.Config) {
			t.Error("mismatch", spew.Sdump(forks[f], mg.Config))
		}

	}
}
