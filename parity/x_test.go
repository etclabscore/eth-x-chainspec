package parity

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"math/big"
// 	"os"
// 	"path/filepath"
// 	"reflect"
// 	"strings"
// 	"testing"

// 	"github.com/davecgh/go-spew/spew"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core"
// 	"github.com/go-test/deep"
// )

// var xreferenceSupportedConfigs = map[string]*core.Genesis{
// 	"foundation.json": core.DefaultGenesisBlock(),
// 	"classic.json":    core.DefaultClassicGenesisBlock(),
// 	"ropsten.json":    core.DefaultTestnetGenesisBlock(),
// 	"mix.json":        core.DefaultMixGenesisBlock(),
// }

// func TestX1(t *testing.T) {
// 	fis, err := ioutil.ReadDir(testChainsJSONDir)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, f := range fis {
// 		fname := filepath.Join(testChainsJSONDir, f.Name())
// 		b, err := ioutil.ReadFile(fname)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		c := &Config{}
// 		err = json.Unmarshal(b, &c)
// 		if err != nil {
// 			t.Fatal(fname, err)
// 		}

// 		mg := c.ToMultiGethGenesis()
// 		if mg == nil {
// 			t.Log("skipping unsupported config", fname)
// 			continue
// 		}

// 		if c.Genesis == nil {
// 			t.Log("config read no genesis")
// 			return
// 		}

// 		if c.Genesis.StateRoot != nil {
// 			mgb := mg.ToBlock(nil)
// 			gotRoot := mgb.Root()
// 			wantRoot := c.Genesis.StateRoot
// 			if gotRoot != *wantRoot {
// 				t.Errorf("%s - got: %x, want: %x", fname, gotRoot, wantRoot)
// 				if f.Name() == "classic.json" {
// 					diff := deep.Equal(mg.Alloc, core.DefaultClassicGenesisBlock().Alloc)
// 					for _, d := range diff {
// 						t.Log(d)
// 					}
// 					for k, v := range mg.Alloc {
// 						ck, ok := core.DefaultClassicGenesisBlock().Alloc[k]
// 						if !ok {
// 							t.Error("missing key A", k, ck, v)
// 						}
// 					}
// 					for k, v := range core.DefaultClassicGenesisBlock().Alloc {
// 						_, ok := mg.Alloc[k]
// 						if !ok {
// 							t.Error("missing key B", k.Hex(), spew.Sdump(v))
// 						}
// 					}
// 				}
// 			}
// 		} else {
// 			// t.Log(fname, "skipping hardcoded stateroot check (DNE)")
// 		}

// 		if f.Name() == "morden.json" {
// 			mgb := mg.ToBlock(nil)
// 			gotRoot := mgb.Root()
// 			wantMordenStateRoot := common.HexToHash("0xf3f4696bbf3b3b07775128eb7a3763279a394e382130f27c21e70233e04946a9")
// 			if gotRoot != wantMordenStateRoot {
// 				t.Errorf("%s - got: %x, want: %x", fname, gotRoot, wantMordenStateRoot)
// 			}
// 		}

// 		wantG, ok := xreferenceSupportedConfigs[f.Name()]
// 		spew.Config.Indent = "\t"
// 		spew.Config.DisableMethods = true
// 		if ok {
// 			// FIXME: WHY IS THIS PASSING?
// 			// The read values should be setting different fields than their corresponding hardcoded equivalent config.
// 			// The read values prefer the FEATURE based fields, while the hardcoded configs still use the hardfork fields.
// 			// So I would expect the DeepEquals checks to say that the struct values are NOT equal.
// 			t.Log("comparing configs read vs hardcoded", f.Name())
// 			if diff := deep.Equal(wantG, mg); diff != nil {
// 				for _, d := range diff {
// 					if !strings.Contains(d, "EIP150Hash") {
// 						t.Error(fname, d)
// 					}
// 				}
// 			} else if !reflect.DeepEqual(wantG, mg) {
// 				t.Error(spew.Sdump(wantG), spew.Sdump(mg))

// 			} else {
// 				// debugging the fixme above
// 				// t.Log(spew.Sdump(wantG), spew.Sdump(mg))
// 			}
// 		}
// 		// also debugging the fixme above
// 		// if f.Name() == "foundation.json" {
// 		// scs := spew.ConfigState{Indent: "\t", DisableMethods: true}
// 		// t.Log(scs.Sdump(params.MainnetChainConfig))
// 		// t.Log(scs.Sdump(mg.Config))
// 		// }
// 	}
// }

// func TestTransitionConfig(t *testing.T) {
// 	fname := filepath.Join(testChainsJSONDir, "transition_test.json")
// 	b, err := ioutil.ReadFile(fname)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	c := &Config{}
// 	err = json.Unmarshal(b, &c)
// 	if err != nil {
// 		t.Fatal(fname, err)
// 	}

// 	mg := c.ToMultiGethGenesis()
// 	if mg == nil {
// 		t.Fatal("skipping unsupported config", fname)
// 	}

// 	if c.Genesis == nil {
// 		t.Fatal("config read no genesis")
// 	}

// 	// mgGenesis := mg.ToBlock(nil)

// 	// hardcode := &params.ChainConfig{
// 	// 	ChainID:        big.NewInt(1),
// 	// 	HomesteadBlock: big.NewInt(0),
// 	// 	EIP150Block:    big.NewInt(0),
// 	// 	EIP155Block:    big.NewInt(0),
// 	// 	EIP158Block:    big.NewInt(0),
// 	// 	ByzantiumBlock: big.NewInt(5),
// 	// }

// }

// func TestBuiltinActivateAt(t *testing.T) {
// 	fname := filepath.Join(testChainsJSONDir, "transition_test.json")
// 	b, err := ioutil.ReadFile(fname)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	c := &Config{}
// 	err = json.Unmarshal(b, &c)
// 	if err != nil {
// 		t.Fatal(fname, err)
// 	}

// 	mg := c.ToMultiGethGenesis()
// 	if mg == nil {
// 		t.Fatal("skipping unsupported config", fname)
// 	}

// 	if c.Genesis == nil {
// 		t.Fatal("config read no genesis")
// 	}
// 	zero := big.NewInt(0)
// 	if mg.Config.IsEIP198F(zero) {
// 		t.Fatal("no")
// 	}
// 	if mg.Config.IsEIP212F(zero) {
// 		t.Fatal("no")
// 	}
// 	if mg.Config.IsEIP213F(zero) {
// 		t.Fatal("no")
// 	}

// 	five := big.NewInt(5)
// 	if !mg.Config.IsEIP198F(five) {
// 		t.Fatal("no")
// 	}
// 	if !mg.Config.IsEIP212F(five) {
// 		t.Fatal("no")
// 	}
// 	if !mg.Config.IsEIP213F(five) {
// 		t.Fatal("no")
// 	}

// 	if !mg.Config.IsEIP161F(five) {
// 		t.Fatal("no")
// 	}

// 	// "EIP158ToByzantiumAt5": {
// 	// 	ChainID:        big.NewInt(1),
// 	// 	HomesteadBlock: big.NewInt(0),
// 	// 	EIP150Block:    big.NewInt(0),
// 	// 	EIP155Block:    big.NewInt(0),
// 	// 	EIP158Block:    big.NewInt(0),
// 	// 	ByzantiumBlock: big.NewInt(5),
// 	// },

// 	// hardcode := &params.ChainConfig{
// 	// 	ChainID:        big.NewInt(1),
// 	// 	HomesteadBlock: big.NewInt(0),
// 	// 	EIP150Block:    big.NewInt(0),
// 	// 	EIP155Block:    big.NewInt(0),
// 	// 	EIP158Block:    big.NewInt(0),
// 	// 	ByzantiumBlock: big.NewInt(5),
// 	// }

// 	// if diff := deep.Equal(hardcode, mg.Config); len(diff) != 0 {
// 	// 	for _, d := range diff {
// 	// 		t.Error(d)
// 	// 	}
// 	// }

// 	// w := c.EngineOpt.ParityConfigEngineEthash.Params.DifficultyBombDelays[xchain.Uint64(5)]
// 	// if *w != xchain.Uint64(3000000) {
// 	// 	t.Fatal("no", spew.Sdump(c.EngineOpt.ParityConfigEngineEthash))
// 	// }

// 	// t.Log()

// 	// if mg.Config.DifficultyBombDelays[big.NewInt(5)] == nil || mg.Config.DifficultyBombDelays[big.NewInt(5)].Cmp(big.NewInt(3000000)) != 0 {
// 	// 	spew.Config.DisableMethods = true
// 	// 	t.Fatal("no", spew.Sdump(mg.Config))
// 	// }
// }

// func TestMultiGethToParityConfig(t *testing.T) {

// 	cases := []struct {
// 		name string
// 		gen  *core.Genesis
// 	}{
// 		{
// 			"classic",
// 			core.DefaultClassicGenesisBlock(),
// 		},
// 	}
// 	for _, c := range cases {
// 		pc := &Config{}
// 		if err := pc.FromMultiGethGenesis("classic-test", c.gen); err != nil {
// 			t.Fatal(err)
// 		}

// 		fname := filepath.Join(testChainsJSONDir+"_out", fmt.Sprintf("multigeth_%s.json", c.name))

// 		b, err := json.MarshalIndent(pc, "", "    ")
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		err = ioutil.WriteFile(fname, b, os.ModePerm)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// }
