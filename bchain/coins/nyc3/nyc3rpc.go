package nyc3

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
    "github.com/juju/errors"
	"encoding/json"
	"github.com/golang/glog"
)

// Nyc3RPC is an interface to JSON-RPC bitcoind service.
type Nyc3RPC struct {
	*btc.BitcoinRPC
}

// getnetworkinfo
type CmdGetNetworkInfo struct {
        Method string `json:"method"`
}
type ResGetNetworkInfo struct {
        Error  *bchain.RPCError `json:"error"`
        Result struct {
                Version         json.Number `json:"version"`
                Subversion      json.Number `json:"subversion"`
                ProtocolVersion json.Number `json:"protocolversion"`
                Timeoffset      float64     `json:"timeoffset"`
                Warnings        string      `json:"warnings"`
        } `json:"result"`
}

// getblockchaininfo
type CmdGetNyc3BlockChainInfo struct {
        Method string `json:"method"`
}
type ResGetBlockChainInfo struct {
        Error  *bchain.RPCError `json:"error"`
        Result struct {
                Chain         string      `json:"chain"`
                Blocks        int         `json:"blocks"`
                Headers       int         `json:"headers"`
                Bestblockhash string      `json:"bestblockhash"`
                SizeOnDisk    int64       `json:"size_on_disk"`
                Warnings      string      `json:"warnings"`
        } `json:"result"`
}

// GetChainInfo returns information about the connected backend
func (b *Nyc3RPC) GetChainInfo() (*bchain.ChainInfo, error) {
        glog.V(1).Info("rpc: getblockchaininfo")

        resCi := ResGetBlockChainInfo{}
        err := b.Call(&CmdGetNyc3BlockChainInfo{Method: "getblockchaininfo"}, &resCi)
        if err != nil {
                return nil, err
        }
        if resCi.Error != nil {
                return nil, resCi.Error
        }

        glog.V(1).Info("rpc: getnetworkinfo")
        resNi := ResGetNetworkInfo{}
        err = b.Call(&CmdGetNetworkInfo{Method: "getnetworkinfo"}, &resNi)
        if err != nil {
                return nil, err
        }
        if resNi.Error != nil {
                return nil, resNi.Error
        }

        // cant store difficulty here as there are two types
        rv := &bchain.ChainInfo{
                Bestblockhash: resCi.Result.Bestblockhash,
                Blocks:        resCi.Result.Blocks,
                Chain:         resCi.Result.Chain,
                Headers:       resCi.Result.Headers,
                SizeOnDisk:    resCi.Result.SizeOnDisk,
                Subversion:    string(resNi.Result.Subversion),
                Timeoffset:    resNi.Result.Timeoffset,
        }
        rv.Version = string(resNi.Result.Version)
        rv.ProtocolVersion = string(resNi.Result.ProtocolVersion)
        if len(resCi.Result.Warnings) > 0 {
                rv.Warnings = resCi.Result.Warnings + " "
        }
        if resCi.Result.Warnings != resNi.Result.Warnings {
                rv.Warnings += resNi.Result.Warnings
        }
        return rv, nil
}


// NewNyc3RPC returns new Nyc3RPC instance.
func NewNyc3RPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &Nyc3RPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV1{}
	s.ChainConfig.SupportsEstimateSmartFee = true

	return s, nil
}

// GetBlock returns block with given hash.
func (b *Nyc3RPC) GetBlock(hash string, height uint32) (*bchain.Block, error) {
	var err error
	if hash == "" && height > 0 {
		hash, err = b.GetBlockHash(height)
		if err != nil {
			return nil, err
		}
	}

	glog.V(1).Info("rpc: getblock (verbosity=1) ", hash)

	res := btc.ResGetBlockThin{}
	req := btc.CmdGetBlock{Method: "getblock"}
	req.Params.BlockHash = hash
	req.Params.Verbosity = 1
	err = b.Call(&req, &res)

	if err != nil {
		return nil, errors.Annotatef(err, "hash %v", hash)
	}
	if res.Error != nil {
		return nil, errors.Annotatef(res.Error, "hash %v", hash)
	}

	txs := make([]bchain.Tx, 0, len(res.Result.Txids))
	for _, txid := range res.Result.Txids {
		tx, err := b.GetTransaction(txid)
		if err != nil {
			if err == bchain.ErrTxNotFound {
				glog.Errorf("rpc: getblock: skipping transaction in block %s due error: %s", hash, err)
				continue
			}
			return nil, err
		}
		txs = append(txs, *tx)
	}
	block := &bchain.Block{
		BlockHeader: res.Result.BlockHeader,
		Txs:         txs,
	}
	return block, nil
}

func isErrBlockNotFound(err *bchain.RPCError) bool {
	return err.Message == "Block not found" ||
		err.Message == "Block height out of range"
}
