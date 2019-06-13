package nyc3

import (
	"blockbook/bchain/coins/btc"
        "github.com/martinboehm/btcd/wire"
        "github.com/martinboehm/btcutil/chaincfg"
)

const (
	MainnetMagic wire.BitcoinNet = 0x1fbeeff1
	TestnetMagic wire.BitcoinNet = 0x1fbeeff1
)

var (
	MainNetParams chaincfg.Params
	TestNetParams chaincfg.Params
)

func init() {
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{53}
	MainNetParams.ScriptHashAddrID = []byte{5}
	MainNetParams.Bech32HRPSegwit = "ny"

	TestNetParams = chaincfg.TestNet3Params
	TestNetParams.Net = TestnetMagic
	TestNetParams.PubKeyHashAddrID = []byte{53}
	TestNetParams.ScriptHashAddrID = []byte{5}
	TestNetParams.Bech32HRPSegwit = "ny"
}

// NYC3Parser handle
type NYC3Parser struct {
	*btc.BitcoinParser
}

// NewNYC3Parser returns new DashParser instance
func NewNYC3Parser(params *chaincfg.Params, c *btc.Configuration) *NYC3Parser {
	return &NYC3Parser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
	}
}

// GetChainParams contains network parameters for the main NYC3 network,
// the regression test NYC3 network, the test NYC3 network and
// the simulation test NYC3 network, in this order
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
