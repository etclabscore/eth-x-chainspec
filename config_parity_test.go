package xchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var testChainsJSONDir = "./parity-chainspecs"

func TestUint64UnmarshalJSON(t *testing.T) {
	ex1 := `"0xC3500"`
	u := new(Uint64)
	err := u.UnmarshalJSON([]byte(ex1))
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONUnmarshaling(t *testing.T) {
	fis, err := ioutil.ReadDir(testChainsJSONDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fis {
		fname := filepath.Join(testChainsJSONDir, f.Name())
		if err := testChainFile(fname); err != nil {
			t.Fatal(err)
		}
	}
}

type chainMarshalCase struct {
	chainFile string
	want      *ParityConfig
}

func mustBlockReward(m map[Uint64]string) *BlockReward {
	br := BlockReward{}
	for k, v := range m {
		wantBR, ok := big.NewInt(0).SetString(v, 16)
		if !ok {
			panic("not ok big string")
		}
		br[k] = (*hexutil.Big)(wantBR)
	}
	return &br
}

func mustBTreeMap(m map[Uint64]Uint64) *BTreeMap {
	bt := BTreeMap{}
	for k, v := range m {
		bt[k] = v
	}
	return &bt
}

var testCases = []chainMarshalCase{
	{
		filepath.Join(testChainsJSONDir, "callisto.json"),
		&ParityConfig{
			Name: "Callisto",
			EngineOpt: ParityConfigEngines{
				ParityConfigEngineEthash: &ParityConfigEngineEthash{
					Params: ParityConfigEngineEthashParams{
						MinimumDifficulty:   Uint64(131072),
						HomesteadTransition: Uint64(0),
						BlockReward: mustBlockReward(
							map[Uint64]string{
								Uint64(0): "16c4abbebea0100000",
							},
						),
						EIP100BTransition: Uint64(20),
						DifficultyBombDelays: mustBTreeMap(map[Uint64]Uint64{
							Uint64(20): Uint64(3000000),
						}),
					},
				},
			},
		},
	},
	{
		filepath.Join(testChainsJSONDir, "foundation.json"),
		&ParityConfig{
			Name: "Ethereum",
			EngineOpt: ParityConfigEngines{
				ParityConfigEngineEthash: &ParityConfigEngineEthash{
					Params: ParityConfigEngineEthashParams{
						MinimumDifficulty:   Uint64(131072),
						HomesteadTransition: Uint64(1150000),
						BlockReward: mustBlockReward(
							map[Uint64]string{
								Uint64(0):       "4563918244f40000",
								Uint64(4370000): "29a2241af62c0000",
								Uint64(7280000): "1bc16d674ec80000",
							},
						),
						EIP100BTransition: Uint64(4370000),
						DifficultyBombDelays: mustBTreeMap(map[Uint64]Uint64{
							Uint64(4370000): Uint64(3000000),
							Uint64(7280000): Uint64(2000000),
						}),
					},
				},
			},
		},
	},
}

func testChainFile(f string) (err error) {
	by, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	p := ParityConfig{}
	err = json.Unmarshal(by, &p)
	if err != nil {
		return fmt.Errorf("%s - %s", f, err)
	}

	defer func() {
		if err != nil {
			fmt.Println(spew.Sdump(p))
		}
	}()

	for _, c := range testCases {
		if c.chainFile == f {
			if c.want.Name != p.Name {
				return fmt.Errorf("%s - got: %v, want: %v", c.chainFile, p.Name, c.want.Name)
			}
			if c.want.EngineOpt.ParityConfigEngineEthash != nil {
				err := assertSameEthashParams(c.chainFile, c.want.EngineOpt.ParityConfigEngineEthash, p.EngineOpt.ParityConfigEngineEthash)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func assertSameEthashParams(chainFile string, p1, p2 *ParityConfigEngineEthash) error {
	if p1.Params.MinimumDifficulty != p1.Params.MinimumDifficulty {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.Params.MinimumDifficulty, p1.Params.MinimumDifficulty)
	}
	if p1.Params.HomesteadTransition != p1.Params.HomesteadTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.Params.HomesteadTransition, p1.Params.HomesteadTransition)
	}
	if p1.Params.EIP100BTransition != p1.Params.EIP100BTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.Params.EIP100BTransition, p1.Params.EIP100BTransition)
	}

	if !reflect.DeepEqual(p1.Params.BlockReward, p1.Params.BlockReward) {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.Params.BlockReward, p1.Params.BlockReward)
	}
	if !reflect.DeepEqual(p1.Params.DifficultyBombDelays, p1.Params.DifficultyBombDelays) {
		if len(*p1.Params.DifficultyBombDelays) > 0 || len(*p1.Params.DifficultyBombDelays) > 0 {

			return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.Params.DifficultyBombDelays, p1.Params.DifficultyBombDelays)
		}
	}
	return nil
}
