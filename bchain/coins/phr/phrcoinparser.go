package phr

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/jakm/btcutil/base58"
	"github.com/jakm/btcutil/chaincfg"
)

const (
	// MainNet represents the main bitcoin network.
	MainPhoreNet wire.BitcoinNet = 0x504852   // PHR
	TestPhoreNet wire.BitcoinNet = 0x54504852 // TPHR
)

var PhoreMainNetParams = chaincfg.Params{
	Name:        "mainPhore",
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
	Checkpoints: []chaincfg.Checkpoint{},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "ph", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID: 0x37, // starts with 1
	ScriptHashAddrID: 0x0d, // starts with 3
	PrivateKeyID:     0xd4, // starts with 5 (uncompressed) or K (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x02, 0x2d, 0x25, 0x33}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x02, 0x21, 0x31, 0x2b}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x800001bc,
}

func init() {
	MainNetParams = PhoreMainNetParams

	TestNetParams = PhoreMainNetParams
	TestNetParams.Name = "testPhore"
	TestNetParams.Net = TestPhoreNet
	TestNetParams.Bech32HRPSegwit = "tph"
}

// PhoreParser handle
type PhoreParser struct {
	*btc.BitcoinParser
	baseparser *bchain.BaseParser
}

// NewGroestlcoinParser returns new PhoreParser instance
func NewPhoreParser(params *chaincfg.Params, c *btc.Configuration) *PhoreParser {
	return &PhoreParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
}

// GetChainParams contains network parameters for the main Groestlcoin network,
// the regression test Groestlcoin network, the test Groestlcoin network and
// the simulation test Groestlcoin network, in this order
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
