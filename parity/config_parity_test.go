package parity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"testing"

	xchain ".."
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var testChainsJSONDir = "./chainspecs"

func TestUint64UnmarshalJSON(t *testing.T) {
	ex1 := `"0xC3500"`
	u := new(xchain.Uint64)
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
	want      *Config
}

func mustBlockReward(m map[xchain.Uint64]string) *xchain.BlockReward {
	br := xchain.BlockReward{}
	for k, v := range m {
		wantBR, ok := big.NewInt(0).SetString(v, 16)
		if !ok {
			panic("not ok big string")
		}
		br[k] = (*hexutil.Big)(wantBR)
	}
	return &br
}

func mustBTreeMap(m map[xchain.Uint64]xchain.Uint64) *xchain.BTreeMap {
	bt := xchain.BTreeMap{}
	for k, v := range m {
		bt[k] = v
	}
	return &bt
}

var testCases = []chainMarshalCase{
	{
		filepath.Join(testChainsJSONDir, "callisto.json"),
		&Config{
			Name: "Callisto",
			EngineOpt: ConfigEngines{
				ParityConfigEngineEthash: &ConfigEngineEthash{
					Params: ConfigEngineEthashParams{
						MinimumDifficulty:   xchain.Uint64(131072),
						HomesteadTransition: xchain.Uint64(0),
						BlockReward: mustBlockReward(
							map[xchain.Uint64]string{
								xchain.Uint64(0): "16c4abbebea0100000",
							},
						),
						EIP100BTransition: xchain.Uint64(20),
						DifficultyBombDelays: mustBTreeMap(map[xchain.Uint64]xchain.Uint64{
							xchain.Uint64(20): xchain.Uint64(3000000),
						}),
					},
				},
			},
			Params: &ConfigParams{
				GasLimitBoundDivisor: xchain.Uint64(1024),
				Registrar: func() *common.Address {
					a := common.HexToAddress("0x0000000000000000000000000000000000000000")
					return &a
				}(),
				AccountStartNonce:     xchain.Uint64(0),
				MaximumExtraDataSize:  xchain.Uint64(32),
				MinGasLimit:           xchain.Uint64(5000),
				NetworkID:             xchain.Uint64(1),
				ChainID:               xchain.Uint64(uint64(0x0334)), // shoulda done 'em all like this; removes 'magic' from conversion from raw json file
				MaxCodeSize:           xchain.Uint64(24576),
				MaxCodeSizeTransition: xchain.Uint64(10),
				EIP150Transition:      xchain.Uint64(0),
				EIP160Transition:      xchain.Uint64(10),
				EIP161abcTransition:   xchain.Uint64(10),
				EIP161dTransition:     xchain.Uint64(10),
				EIP155Transition:      xchain.Uint64(10),
				EIP140Transition:      xchain.Uint64(20),
				EIP211Transition:      xchain.Uint64(20),
				EIP214Transition:      xchain.Uint64(20),
				EIP658Transition:      xchain.Uint64(20),
			},
		},
	},
	{
		filepath.Join(testChainsJSONDir, "foundation.json"),
		&Config{
			Name: "Ethereum",
			EngineOpt: ConfigEngines{
				ParityConfigEngineEthash: &ConfigEngineEthash{
					Params: ConfigEngineEthashParams{
						MinimumDifficulty:   xchain.Uint64(131072),
						HomesteadTransition: xchain.Uint64(1150000),
						BlockReward: mustBlockReward(
							map[xchain.Uint64]string{
								xchain.Uint64(0):       "4563918244f40000",
								xchain.Uint64(4370000): "29a2241af62c0000",
								xchain.Uint64(7280000): "1bc16d674ec80000",
							},
						),
						EIP100BTransition: xchain.Uint64(4370000),
						DifficultyBombDelays: mustBTreeMap(map[xchain.Uint64]xchain.Uint64{
							xchain.Uint64(4370000): xchain.Uint64(3000000),
							xchain.Uint64(7280000): xchain.Uint64(2000000),
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

	p := Config{}
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
				err := assertEthashParams(c.chainFile, c.want.EngineOpt.ParityConfigEngineEthash, p.EngineOpt.ParityConfigEngineEthash)
				if err != nil {
					return err
				}
			}
			if c.want.Params != nil {
				err := assertParams(c.chainFile, c.want.Params, p.Params)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func assertParams(chainFile string, p1, p2 *ConfigParams) error {
	if p1.GasLimitBoundDivisor != p2.GasLimitBoundDivisor {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.GasLimitBoundDivisor, p2.GasLimitBoundDivisor)
	}
	if !reflect.DeepEqual(p1.Registrar, p2.Registrar) {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.Registrar, p2.Registrar)
	}
	if p1.AccountStartNonce != p2.AccountStartNonce {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.AccountStartNonce, p2.AccountStartNonce)
	}
	if p1.MaximumExtraDataSize != p2.MaximumExtraDataSize {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MaximumExtraDataSize, p2.MaximumExtraDataSize)
	}
	if p1.MinGasLimit != p2.MinGasLimit {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MinGasLimit, p2.MinGasLimit)
	}
	if p1.NetworkID != p2.NetworkID {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.NetworkID, p2.NetworkID)
	}
	if p1.ChainID != p2.ChainID {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.ChainID, p2.ChainID)
	}
	if p1.MaxCodeSize != p2.MaxCodeSize {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MaxCodeSize, p2.MaxCodeSize)
	}
	if p1.MaxCodeSizeTransition != p2.MaxCodeSizeTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MaxCodeSizeTransition, p2.MaxCodeSizeTransition)
	}
	if p1.EIP150Transition != p2.EIP150Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP150Transition, p2.EIP150Transition)
	}
	if p1.EIP160Transition != p2.EIP160Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP160Transition, p2.EIP160Transition)
	}
	if p1.EIP161abcTransition != p2.EIP161abcTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP161abcTransition, p2.EIP161abcTransition)
	}
	if p1.EIP161dTransition != p2.EIP161dTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP161dTransition, p2.EIP161dTransition)
	}
	if p1.EIP155Transition != p2.EIP155Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP155Transition, p2.EIP155Transition)
	}
	if p1.EIP140Transition != p2.EIP140Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP140Transition, p2.EIP140Transition)
	}
	if p1.EIP211Transition != p2.EIP211Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP211Transition, p2.EIP211Transition)
	}
	if p1.EIP214Transition != p2.EIP214Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP214Transition, p2.EIP214Transition)
	}
	if p1.EIP658Transition != p2.EIP658Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP658Transition, p2.EIP658Transition)
	}
	return nil
}

func assertEthashParams(chainFile string, p1, p2 *ConfigEngineEthash) error {
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
