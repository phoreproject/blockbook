package phr

import (
	"blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/bchain"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/jakm/btcutil/chaincfg"
)

const (
	// MainNet represents the main bitcoin network.
	MainPhoreNet wire.BitcoinNet = 0x504852   // PHR
	TestPhoreNet wire.BitcoinNet = 0x545048   // TP
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

func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}

var PhoreMainNetParams = chaincfg.Params{
	Name:        "main",
	Net:         MainPhoreNet,
	DefaultPort: "11771",
	DNSSeeds: []chaincfg.DNSSeed{
		{"dns0.phore.io", true},
		{"phore.seed.rho.industries", true},
	},

	// Chain parameters
	GenesisBlock:     nil, // unused
	GenesisHash:      nil, // unused
	PowLimit:         mainPowLimit,
	PowLimitBits:     0x207fffff,
	BIP0034Height:    0, // unused
	BIP0065Height:    0, // unused
	BIP0066Height:    0, // unused
	CoinbaseMaturity: 50,
	TargetTimespan:   time.Minute, // 1 minute
	//PoSTargetTimespan:        time.Minute * 40,
	TargetTimePerBlock:       time.Minute, // 1 minutes
	RetargetAdjustmentFactor: 4,           // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        true,
	//MasternodeDriftCount:     20,
	//LastPoWBlock:             200,
	//ZerocoinStartHeight:      89993,
	//ZerocoinLastOldParams:    99999999,
	//StakeMinimumAge:          time.Hour * 3,
	//ModifierV2StartBlock:     433160,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []chaincfg.Checkpoint{
		{Height: 500000, Hash: newHashFromStr("efaab512d538c5de130a5e11459da65264699f552d69d7a960b0b8effabd97cb")},
		{Height: 1000000, Hash: newHashFromStr("a88a1ad20733dda703a73c12916c0b77e6242a71b7f88ad827aee8bd9264c652")},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "ph", // always bc for main net

	AddressMagicLen: 1,

	// Address encoding magics
	PubKeyHashAddrID: []byte{0x37}, // starts with 1
	ScriptHashAddrID: []byte{0x0d}, // starts with 3
	PrivateKeyID:     []byte{0xd4}, // starts with 5 (uncompressed) or K (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x02, 0x2d, 0x25, 0x33}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x02, 0x21, 0x31, 0x2b}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x800001bc,
}

var PhoreTestNetParams = chaincfg.Params{
	Name:        "testnetPhore",
	Net:         TestPhoreNet,
	DefaultPort: "11773",
	DNSSeeds: []chaincfg.DNSSeed{
	},

	// Chain parameters
	GenesisBlock:     nil, // unused
	GenesisHash:      nil, // unused
	PowLimit:         mainPowLimit,
	PowLimitBits:     0x207fffff,
	BIP0034Height:    0, // unused
	BIP0065Height:    0, // unused
	BIP0066Height:    0, // unused
	CoinbaseMaturity: 50,
	TargetTimespan:   time.Minute, // 1 minute
	TargetTimePerBlock:       time.Minute, // 1 minutes
	RetargetAdjustmentFactor: 4,           // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        true,
	Checkpoints: []chaincfg.Checkpoint{},
	RelayNonStdTxs: false,
	Bech32HRPSegwit: "tp",

	AddressMagicLen: 1,

	PubKeyHashAddrID: []byte{0x8B}, // starts with x or y
	ScriptHashAddrID: []byte{0x13}, // starts with 8 or 9
	PrivateKeyID:     []byte{0xEF}, // starts with '9' or 'c' (Bitcoin defaults)

	HDPrivateKeyID: [4]byte{0x3a, 0x80, 0x61, 0xa0},
	HDPublicKeyID:  [4]byte{0x3a, 0x80, 0x58, 0x37},

	HDCoinType: 0x80000001,
}


func init() {
	MainNetParams = PhoreMainNetParams
	TestNetParams = PhoreTestNetParams
}

// PhoreParser handle
type PhoreParser struct {
	*btc.BitcoinParser
	baseparser *bchain.BaseParser
}

// NewPhoreParser returns new PhoreParser instance
func NewPhoreParser(params *chaincfg.Params, c *btc.Configuration) *PhoreParser {
	return &PhoreParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
}

// GetChainParams contains network parameters for the main Phore network,
// and the test Phore network
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
func (p *PhoreParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *PhoreParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
