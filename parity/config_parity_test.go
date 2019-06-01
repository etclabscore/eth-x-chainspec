package parity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	xchain ".."
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-test/deep"
)

var testChainsJSONDir = "./chainspecs"

func TestUint64UnmarshalJSON(t *testing.T) {
	ex1 := `"0xC3500"`
	u := new(xchain.Uint64)
	err := u.UnmarshalJSON([]byte(ex1))
	if err != nil {
		t.Fatal(err)
	}

	ex2 := `"0x1"`
	u = new(xchain.Uint64)
	err = u.UnmarshalJSON([]byte(ex2))
	if err != nil {
		t.Fatal(err)
	}
	if *u != 1 {
		t.Fatalf("got: %v, want: %v", *u, 1)
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

func mustBTreeMap(m map[xchain.Uint64]*xchain.Uint64) *xchain.BTreeMap {
	bt := xchain.BTreeMap{}
	for k, v := range m {
		bt[k] = v
	}
	return &bt
}

func xchainUint64(u uint64) *xchain.Uint64 {
	x := xchain.Uint64(u)
	return &x
}

var testCases = []chainMarshalCase{
	{
		filepath.Join(testChainsJSONDir, "callisto.json"),
		&Config{
			Name: "Callisto",
			EngineOpt: ConfigEngines{
				ParityConfigEngineEthash: &ConfigEngineEthash{
					Params: ConfigEngineEthashParams{
						MinimumDifficulty:   xchainUint64(131072),
						HomesteadTransition: xchainUint64(0),
						BlockReward: mustBlockReward(
							map[xchain.Uint64]string{
								*xchainUint64(0): "16c4abbebea0100000",
							},
						),
						EIP100BTransition: xchainUint64(20),
						DifficultyBombDelays: mustBTreeMap(map[xchain.Uint64]*xchain.Uint64{
							*xchainUint64(20): xchainUint64(3000000),
						}),
					},
				},
			},
			Params: &ConfigParams{
				GasLimitBoundDivisor: xchainUint64(uint64(0x0400)),
				Registrar: func() *common.Address {
					a := common.HexToAddress("0x0000000000000000000000000000000000000000")
					return &a
				}(),
				AccountStartNonce:     xchainUint64(0),
				MaximumExtraDataSize:  xchainUint64(32),
				MinGasLimit:           xchainUint64(5000),
				NetworkID:             xchainUint64(1),
				ChainID:               xchainUint64(uint64(0x0334)), // shoulda done 'em all like this; removes 'magic' from conversion from raw json file
				MaxCodeSize:           xchainUint64(24576),
				MaxCodeSizeTransition: xchainUint64(10),
				EIP150Transition:      xchainUint64(0),
				EIP160Transition:      xchainUint64(10),
				EIP161abcTransition:   xchainUint64(10),
				EIP161dTransition:     xchainUint64(10),
				EIP155Transition:      xchainUint64(10),
				EIP140Transition:      xchainUint64(20),
				EIP211Transition:      xchainUint64(20),
				EIP214Transition:      xchainUint64(20),
				EIP658Transition:      xchainUint64(20),
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
						MinimumDifficulty:   xchainUint64(131072),
						HomesteadTransition: xchainUint64(1150000),
						BlockReward: mustBlockReward(
							map[xchain.Uint64]string{
								*xchainUint64(0):       "4563918244f40000",
								*xchainUint64(4370000): "29a2241af62c0000",
								*xchainUint64(7280000): "1bc16d674ec80000",
							},
						),
						EIP100BTransition: xchainUint64(4370000),
						DifficultyBombDelays: mustBTreeMap(map[xchain.Uint64]*xchain.Uint64{
							*xchainUint64(4370000): xchainUint64(3000000),
							*xchainUint64(7280000): xchainUint64(2000000),
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
	if *p1.GasLimitBoundDivisor != *p2.GasLimitBoundDivisor {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.GasLimitBoundDivisor, p2.GasLimitBoundDivisor)
	}
	if !reflect.DeepEqual(p1.Registrar, p2.Registrar) {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.Registrar, p2.Registrar)
	}
	if *p1.AccountStartNonce != *p2.AccountStartNonce {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.AccountStartNonce, p2.AccountStartNonce)
	}
	if *p1.MaximumExtraDataSize != *p2.MaximumExtraDataSize {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MaximumExtraDataSize, p2.MaximumExtraDataSize)
	}
	if *p1.MinGasLimit != *p2.MinGasLimit {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MinGasLimit, p2.MinGasLimit)
	}
	if *p1.NetworkID != *p2.NetworkID {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.NetworkID, p2.NetworkID)
	}
	if *p1.ChainID != *p2.ChainID {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.ChainID, p2.ChainID)
	}
	if *p1.MaxCodeSize != *p2.MaxCodeSize {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MaxCodeSize, p2.MaxCodeSize)
	}
	if *p1.MaxCodeSizeTransition != *p2.MaxCodeSizeTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.MaxCodeSizeTransition, p2.MaxCodeSizeTransition)
	}
	if *p1.EIP150Transition != *p2.EIP150Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP150Transition, p2.EIP150Transition)
	}
	if *p1.EIP160Transition != *p2.EIP160Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP160Transition, p2.EIP160Transition)
	}
	if *p1.EIP161abcTransition != *p2.EIP161abcTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP161abcTransition, p2.EIP161abcTransition)
	}
	if *p1.EIP161dTransition != *p2.EIP161dTransition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP161dTransition, p2.EIP161dTransition)
	}
	if *p1.EIP155Transition != *p2.EIP155Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP155Transition, p2.EIP155Transition)
	}
	if *p1.EIP140Transition != *p2.EIP140Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP140Transition, p2.EIP140Transition)
	}
	if *p1.EIP211Transition != *p2.EIP211Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP211Transition, p2.EIP211Transition)
	}
	if *p1.EIP214Transition != *p2.EIP214Transition {
		return fmt.Errorf("%s - got: %v, want: %v", chainFile, p1.EIP214Transition, p2.EIP214Transition)
	}
	if *p1.EIP658Transition != *p2.EIP658Transition {
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

func TestJSONMarshallUint64(t *testing.T) {
	x := xchain.Uint64(0x42)
	_, err := json.Marshal(&x)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONMarshaling(t *testing.T) {
	by, err := ioutil.ReadFile(filepath.Join(testChainsJSONDir, "classic.json"))
	if err != nil {
		t.Fatal(err)
	}

	p := Config{}
	err = json.Unmarshal(by, &p)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err != nil {
			fmt.Println(spew.Sdump(p))
		}
	}()

	out, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	p2 := Config{}
	err = json.Unmarshal(out, &p2)
	if err != nil {
		t.Fatal(err)
	}

	// Instead of using normal reflect.DeepEqual, this package does the same thing
	// but conveniently shows the diffs between the structs, if any.
	// This was used for the debugging noted below.
	if diff := deep.Equal(p, p2); diff != nil {
		// This debugging was added because if the maps for BTreeMap and BlockReward use pointers as keys,
		// then the test fails because the addresses are not equal.
		// Backwards fitting? Maybe.
		// #A13:sketchybuttestspass
		b, err := json.MarshalIndent(p.EngineOpt.ParityConfigEngineEthash.Params.BlockReward, "", "    ")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
		t.Log(string(out))
		t.Fatal(diff)
	}
}

func TestChainSpecJSONReproducability(t *testing.T) {
	outTestChainsJSONDir := testChainsJSONDir + "_out"
	if err := os.MkdirAll(outTestChainsJSONDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	ffs, err := ioutil.ReadDir(testChainsJSONDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range ffs {
		b, err := ioutil.ReadFile(filepath.Join(testChainsJSONDir, f.Name()))
		if err != nil {
			t.Fatal(err)
		}

		p := Config{}
		err = json.Unmarshal(b, &p)
		if err != nil {
			t.Fatal(spew.Sdump(p), f.Name(), err)
		}

		outb, err := json.MarshalIndent(&p, "", "     ")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(filepath.Join(outTestChainsJSONDir, f.Name()), outb, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
}
