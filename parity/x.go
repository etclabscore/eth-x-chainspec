package parity

import (
	"fmt"
	"math/big"
	"strings"

	xchain ".."

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// ToMultiGethGenesis converts a Parity chainspec to the corresponding MultiGeth datastructure.
// Note that the return value 'core.Genesis' includes the respective 'params.ChainConfig' values.
func (c *Config) ToMultiGethGenesis() *core.Genesis {
	mgc := &params.ChainConfig{}
	if pars := c.Params; pars != nil {
		// FIXME
		if pars.EIP161abcTransition.Uint64() != pars.EIP161dTransition.Uint64() {
			panic("not supported")
		}
		// unsupportedValuesMust := map[interface{}]interface{}{
		// 	pars.AccountStartNonce:                       uint64(0),
		// 	pars.MaximumExtraDataSize:                    uint64(32),
		// 	pars.MinGasLimit:                             uint64(5000),
		// 	pars.SubProtocolName:                         "",
		// 	pars.EIP98Transition:                         nil,
		// 	pars.ValidateChainIDTransition:               nil,
		// 	pars.ValidateChainReceiptsTransition:         nil,
		// 	pars.DustProtectionTransition:                nil,
		// 	pars.NonceCapIncrement:                       nil,
		// 	pars.RemoveDustContracts:                     false,
		// 	pars.EIP210Transition:                        nil,
		// 	pars.EIP210ContractAddress:                   nil,
		// 	pars.EIP210ContractCode:                      nil,
		// 	pars.ApplyReward:                             false,
		// 	pars.TransactionPermissionContract:           nil,
		// 	pars.TransactionPermissionContractTransition: nil,
		// 	pars.WASMActivationTransition:                nil,
		// 	pars.KIP4Transition:                          nil,
		// 	pars.KIP6Transition:                          nil,
		// 	// TODO...
		// }
		// i := -1
		// for k, v := range unsupportedValuesMust {
		// 	i++
		// 	if v == nil && k == nil {
		// 		continue
		// 	}
		// 	if v != nil && !reflect.DeepEqual(k, v) {
		// 		panic(fmt.Sprintf("%d: %v != %v - unsupported configuration value", i, k, v))
		// 	}
		// }

		mgc.ChainID = pars.ChainID.Big()
		if mgc.ChainID == nil && pars.NetworkID != nil {
			mgc.ChainID = pars.NetworkID.Big() // Default according to Parity documentation https://wiki.parity.io/Chain-specification.html
		}

		mgc.EIP150Block = pars.EIP150Transition.Big()
		// mgc.EIP150Hash // TODO? CHT?

		mgc.EIP155Block = pars.EIP155Transition.Big()
		mgc.EIP160FBlock = pars.EIP160Transition.Big()
		mgc.EIP170FBlock = pars.MaxCodeSizeTransition.Big()
		if mgc.EIP170FBlock != nil && uint64(*pars.MaxCodeSize) != uint64(24576) {
			panic(fmt.Sprintf("%v != %v - unsupported configuration value", *pars.MaxCodeSize, 24576))
		}

		mgc.EIP161FBlock = pars.EIP161abcTransition.Big()

		mgc.EIP140FBlock = pars.EIP140Transition.Big()
		mgc.EIP145FBlock = pars.EIP145Transition.Big()
		mgc.EIP211FBlock = pars.EIP211Transition.Big()
		mgc.EIP214FBlock = pars.EIP214Transition.Big()
		mgc.EIP658FBlock = pars.EIP658Transition.Big()
		mgc.EIP1014FBlock = pars.EIP1014Transition.Big()
		mgc.EIP1052FBlock = pars.EIP1052Transition.Big()
		mgc.EIP1283FBlock = pars.EIP1283Transition.Big()
		mgc.PetersburgBlock = pars.EIP1283DisableTransition.Big()

		if pars.ForkBlock != nil && pars.ForkCanonHash != nil {
			if (uint64(*pars.ForkBlock) == params.MainnetChainConfig.DAOForkBlock.Uint64() && *pars.ForkCanonHash == common.HexToHash("0x4985f5ca3d2afbec36529aa96f74de3cc10a2a4a6c44f2157a57d2c6059a11bb")) || (uint64(*pars.ForkBlock) == params.TestnetChainConfig.DAOForkBlock.Uint64() && *pars.ForkCanonHash == common.HexToHash("0x3e12d5c0f8d63fbc5831cc7f7273bd824fa4d0a9a4102d65d99a7ea5604abc00")) {

				mgc.DAOForkBlock = new(big.Int).SetUint64(pars.ForkBlock.Uint64())
				mgc.DAOForkSupport = true
			}
			if uint64(*pars.ForkBlock) == uint64(1920000) && *pars.ForkCanonHash == common.HexToHash("0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f") {
				mgc.DAOForkBlock = new(big.Int).SetUint64(pars.ForkBlock.Uint64())
			}
		}
	}

	if ethc := c.EngineOpt.ParityConfigEngineEthash; ethc != nil {

		pars := ethc.Params

		mgc.Ethash = &params.EthashConfig{}

		mgc.HomesteadBlock = pars.HomesteadTransition.Big()
		mgc.EIP100FBlock = pars.EIP100BTransition.Big()
		mgc.DisposalBlock = pars.BombDefuseTransition.Big()
		mgc.ECIP1010PauseBlock = pars.Ecip1010PauseTransition.Big()
		mgc.ECIP1010Length = func() *big.Int {
			if pars.Ecip1010ContinueTransition != nil {
				return new(big.Int).Sub(pars.Ecip1010ContinueTransition.Big(), pars.Ecip1010PauseTransition.Big())
			} else if pars.Ecip1010PauseTransition == nil && pars.Ecip1010ContinueTransition == nil {
				return nil
			}
			return big.NewInt(0)
		}()
		mgc.ECIP1017EraRounds = pars.Ecip1017EraRounds.Big()

	} else if ethc := c.EngineOpt.ParityConfigEngineClique; ethc != nil {

		pars := ethc.Params

		mgc.Clique = &params.CliqueConfig{
			Period: pars.Period,
			Epoch:  pars.Epoch,
		}

	} else {
		return nil
	}
	mgg := &core.Genesis{
		Config: mgc,
	}
	if c.Genesis != nil {
		seal := c.Genesis.Seal.Ethereum

		mgg.Nonce = seal.Nonce.Uint64()
		mgg.Mixhash = seal.MixHash
		mgg.Timestamp = c.Genesis.Timestamp.Uint64()
		mgg.GasLimit = c.Genesis.GasLimit.Uint64()
		mgg.GasUsed = c.Genesis.GasUsed.Uint64()
		mgg.Difficulty = c.Genesis.Difficulty.Big()
		mgg.Coinbase = *c.Genesis.Author
		mgg.ParentHash = *c.Genesis.ParentHash
		mgg.ExtraData = c.Genesis.ExtraData
	}
	if c.Accounts != nil {
		mgg.Alloc = core.GenesisAlloc{}

	accountsloop:
		for k, v := range c.Accounts {
			bal, ok := xchain.ParseBig256(v.Balance)
			if !ok {
				panic("error setting genesis account balance")
			}
			var nonce uint64
			if v.Nonce != nil {
				nonce = uint64(*v.Nonce)
			}

			addr := common.HexToAddress(strings.ToLower(k))
			if _, ok := vm.PrecompiledContractsForConfig(params.AllEthashProtocolChanges, big.NewInt(0))[addr]; ok && bal.Sign() < 1 {
				continue accountsloop
			}

			mgg.Alloc[addr] = core.GenesisAccount{
				Nonce:   nonce,
				Balance: bal,
				Code:    v.Code,
				Storage: v.Storage,
			}
		}
	}
	return mgg
}
