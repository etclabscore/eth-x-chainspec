package parity

import (
	xchain ".."
	"github.com/ethereum/go-ethereum/common"
)

// ParityConfig is the data structure for Parity-Ethereum's chain configuration.
type ParityConfig struct {
	Name      string              `json:"name"`
	DataDir   string              `json:"dataDir"`
	EngineOpt ParityConfigEngines `json:"engine"`
	Params    *ParityConfigParams `json:"params"`
}

type ParityConfigEngines struct {
	ParityConfigEngineEthash         *ParityConfigEngineEthash         `json:"Ethash,omitempty"`
	ParityConfigEngineInstantSeal    *ParityConfigEngineInstantSeal    `json:"instantSeal,omitempty"`
	ParityConfigEngineClique         *ParityConfigEngineClique         `json:"Clique,omitempty"`
	ParityConfigEngineAuthorityRound *ParityConfigEngineAuthorityRound `json:"authorityRound,omitempty"`
}

// ParityConfigEngine is the data structure for a consensus engine.
type ParityConfigEngineEthash struct {
	Params ParityConfigEngineEthashParams `json:"params"`
}

// ParityConfigEngineParamsEthash is the data structure for the Ethash consensus engine parameters.
type ParityConfigEngineEthashParams struct {
	MinimumDifficulty                    xchain.Uint64 `json:"minimumDifficulty,omitempty"`
	DifficultyBoundDivisor               xchain.Uint64 `json:"difficultyBoundDivisor,omitempty"`
	DifficultyIncrementDivisor           xchain.Uint64 `json:"difficultyIncrementDivisor,omitempty"`
	MetropolisDifficultyIncrementDivisor xchain.Uint64 `json:"metropolisDifficultyIncrementDivisor,omitempty"`
	DurationLimit                        xchain.Uint64 `json:"durationLimit,omitempty"`

	HomesteadTransition           xchain.Uint64       `json:"homesteadTransition,omitempty"`
	BlockReward                   *xchain.BlockReward `json:"blockReward,omitempty"`
	BlockRewardContractTransition xchain.Uint64       `json:"blockRewardContractTransition,omitempty"`
	BlockRewardContractAddress    *common.Address     `json:"blockRewardContractAddress,omitempty"`
	BlockRewardContractCode       []byte              `json:"blockRewardContractCode,omitempty"`

	DaoHardforkTransition  xchain.Uint64    `json:"daoHardforkTransition,omitempty"`
	DaoHardforkBeneficiary *common.Address  `json:"daoHardforkBeneficiary,omitempty"`
	DaoHardforkAccounts    []common.Address `json:"daoHardforkAccounts,omitempty"`

	DifficultyHardforkTransition   xchain.Uint64 `json:"difficultyHardforkTransition,omitempty"`
	DifficultyHardforkBoundDivisor xchain.Uint64 `json:"difficultyHardforkBoundDivisor,omitempty"`
	BombDefuseTransition           xchain.Uint64 `json:"bombDefuseTransition,omitempty"`

	EIP100BTransition xchain.Uint64 `json:"eip100bTransition,omitempty"`

	Ecip1010PauseTransition    xchain.Uint64 `json:"ecip1010PauseTransition,omitempty"`
	Ecip1010ContinueTransition xchain.Uint64 `json:"ecip1010ContinueTransition,omitempty"`

	Ecip1017EraRounds xchain.Uint64 `json:"ecip1017EraRounds,omitempty"`

	DifficultyBombDelays *xchain.BTreeMap `json:"difficultyBombDelays,omitempty"`

	EXPIP2Transition    xchain.Uint64 `json:"expip2Transition,omitempty"`
	EXPIP2DurationLimit xchain.Uint64 `json:"expip2DurationLimit,omitempty"`
	ProgPowTransition   xchain.Uint64 `json:"progPowTransition,omitempty"`
}

type ParityConfigEngineInstantSeal struct {
	Params ParityConfigEngineInstantSealParams `json:"params"`
}

type ParityConfigEngineInstantSealParams struct {
	MillisecondTimestamp bool `json:"millisecondTimestamp,omitempty"`
}

type ParityConfigEngineClique struct {
	Params ParityConfigEngineCliqueParams `json:"params"`
}

type ParityConfigEngineCliqueParams struct {
	Period xchain.Uint64 `json:"period,omitempty"`
	Epoch  xchain.Uint64 `json:"epoch,omitempty"`
}

type ParityConfigEngineAuthorityRound struct {
	Params ParityConfigEngineAuthorityRoundParams `json:"params"`
}

type ParityConfigEngineAuthorityRoundParams struct {
}

type ParityConfigParams struct {
}