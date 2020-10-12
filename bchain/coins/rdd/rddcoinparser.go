package rdd

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/jakm/btcutil/chaincfg"
)

const (
	// MainNet represents the main bitcoin network.
	MainReddNet wire.BitcoinNet = 0x504852 // PHR
	TestReddNet wire.BitcoinNet = 0x545048 // TP
)

var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// mainPowLimit is the highest proof of work value a Bitcoin block can
	// have for the main network.  It is the value 2^224 - 1.
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)

	MainNetParams chaincfg.Params
	TestNetParams chaincfg.Params
)

var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xb8, 0x68, 0xe0, 0xd9, 0x5a, 0x3c, 0x3c, 0x0e,
	0x0d, 0xad, 0xc6, 0x7e, 0xe5, 0x87, 0xaa, 0xf9,
	0xdc, 0x8a, 0xcb, 0xf9, 0x9e, 0x3b, 0x4b, 0x31,
	0x10, 0xfa, 0xd4, 0xeb, 0x74, 0xc1, 0xde, 0xcc,
})

func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}

var ReddMainNetParams = chaincfg.Params{
	Name:        "mainRedd",
	Net:         MainReddNet,
	DefaultPort: "45444",
	DNSSeeds: []chaincfg.DNSSeed{
		{"seed.reddcoin.com", true},
		{"dnsseed01.redd.ink", true},
		{"dnsseed02.redd.ink", true},
		{"dnsseed03.redd.ink", true},
	},

	// Chain parameters
	GenesisBlock:             nil, // not required
	GenesisHash:              &genesisHash,
	PowLimit:                 mainPowLimit,
	PowLimitBits:             0x1d00ffff,
	BIP0034Height:            227931,
	BIP0065Height:            388381,
	BIP0066Height:            363725,
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24, // 24 hours
	TargetTimePerBlock:       time.Minute,    // 1 minute
	RetargetAdjustmentFactor: 4,              // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []chaincfg.Checkpoint{
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "bc", // always bc for main net

	AddressMagicLen: 1,

	// Address encoding magics
	PubKeyHashAddrID: []byte{0x3D}, // starts with 61
	ScriptHashAddrID: []byte{0x05}, // starts with 5
	PrivateKeyID:     []byte{0xBD},
	WitnessPubKeyHashAddrID: nil,
	WitnessScriptHashAddrID: nil,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xCF}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x80000004,
}

var ReddTestNetParams = chaincfg.Params{
	Name:        "testnetRedd",
	Net:         TestReddNet,
	DefaultPort: "55444",
	DNSSeeds: []chaincfg.DNSSeed{
		{"testnet-seed.reddcoin.com", true},
		{"testnet-dnsseed.redd.ink", true},
	},

	// Chain parameters
	GenesisBlock:             nil, // unused
	GenesisHash:              nil, // unused
	PowLimit:                 mainPowLimit,
	PowLimitBits:             0x207fffff,
	BIP0034Height:            0, // unused
	BIP0065Height:            0, // unused
	BIP0066Height:            0, // unused
	CoinbaseMaturity:         50,
	TargetTimespan:           time.Minute, // 1 minute
	TargetTimePerBlock:       time.Minute, // 1 minutes
	RetargetAdjustmentFactor: 4,           // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        true,
	Checkpoints:              []chaincfg.Checkpoint{},
	RelayNonStdTxs:           false,
	Bech32HRPSegwit:          "bc",

	AddressMagicLen: 1,

	// Address encoding magics
	PubKeyHashAddrID: []byte{0x8B}, // starts with x or y
	ScriptHashAddrID: []byte{0x13}, // starts with 8 or 9
	PrivateKeyID:     []byte{0xEF}, // starts with '9' or 'c' (Bitcoin defaults)
	WitnessPubKeyHashAddrID: nil,
	WitnessScriptHashAddrID: nil,

	HDPrivateKeyID: [4]byte{0x3a, 0x80, 0x61, 0xa0},
	HDPublicKeyID:  [4]byte{0x3a, 0x80, 0x58, 0x37},

	HDCoinType: 0x80000001,
}


func init() {
	MainNetParams = ReddMainNetParams
	TestNetParams = ReddTestNetParams
}

// ReddParser handle
type ReddParser struct {
	*btc.BitcoinParser
	baseparser *bchain.BaseParser
}

// NewReddParser returns new ReddParser instance
func NewReddParser(params *chaincfg.Params, c *btc.Configuration) *ReddParser {
	return &ReddParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
}

// GetChainParams contains network parameters for the main ReedCoin network,
// and the test ReddCoin network
func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&MainNetParams) {
		err := chaincfg.Register(&MainNetParams)
		if err == nil {
			err = chaincfg.Register(&TestNetParams)
		}
		if err != nil {
			panic(err)
		}
	}
	switch chain {
	case "test":
		return &TestNetParams
	default:
		return &MainNetParams
	}
}

// PackTx packs transaction to byte array using protobuf
func (p *ReddParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *ReddParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
