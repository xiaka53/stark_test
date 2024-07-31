package chain

import (
	"context"
	"encoding/hex"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"log"
	"strconv"
)

var client *rpc.Provider

func init() {
	// c, e := rpc.NewProvider("https://starknet-mainnet.blastapi.io")
	c, e := rpc.NewProvider("https://starknet-mainnet.public.blastapi.io")
	if e != nil {
		log.Fatalf("Failed to connect to starkNet:%v", e)
	}
	client = c
}

func getBlock() (uint64, *rpc.RPCError) {
	return client.BlockNumber(context.Background())
}

func getBlockWithTxs(blockNum uint64) (any, *rpc.RPCError) {
	blockId := rpc.BlockID{Number: &blockNum}
	return client.BlockWithTxs(context.Background(), blockId)
}

func call(contractAddress, function string, callData []*felt.Felt) ([]*felt.Felt, *rpc.RPCError) {
	contract, err := utils.HexToFelt(contractAddress)
	if err != nil {
		return nil, rpc.Err(500, err.Error())
	}
	request := rpc.FunctionCall{
		ContractAddress:    contract,
		EntryPointSelector: utils.GetSelectorFromNameFelt(function),
		Calldata:           callData,
	}
	return client.Call(context.Background(), request, rpc.WithBlockTag("latest"))
}

func getTokenInfo(contractAddress string) (name, symbol string, decimals int) {
	if data, err := call(contractAddress, "name", nil); err == nil {
		if len(data) == 1 {
			_name, e := hex.DecodeString(data[0].String()[2:])
			if e == nil {
				name = string(_name)
			}
		} else {
			_name, e := hex.DecodeString(data[1].String()[2:])
			if e == nil {
				name = string(_name)
			}
		}
	}
	if data, err := call(contractAddress, "symbol", nil); err == nil {
		if len(data) == 1 {
			_symbol, e := hex.DecodeString(data[0].String()[2:])
			if e == nil {
				symbol = string(_symbol)
			}
		} else {
			_symbol, e := hex.DecodeString(data[1].String()[2:])
			if e == nil {
				symbol = string(_symbol)
			}
		}
	}
	if data, err := call(contractAddress, "decimals", nil); err == nil {
		decimals, _ = strconv.Atoi(data[0].Text(10))
	}
	return
}

func isNft721(contractAddress string) bool {
	interfaceID := new(felt.Felt).SetUint64(0x80ac58cd)
	if data, err := call(contractAddress, "supportsInterface", []*felt.Felt{interfaceID}); err == nil {
		if data[0].String() == "0x1" {
			return true
		}
	}
	return false
}

func isNft1155(contractAddress string) bool {
	interfaceID := new(felt.Felt).SetUint64(0xd9b67a26)
	if data, err := call(contractAddress, "supportsInterface", []*felt.Felt{interfaceID}); err == nil {
		if data[0].String() == "0x1" {
			return true
		}
	}
	return false
}
