package chain

import (
	"context"
	"fmt"
	"github.com/NethermindEth/starknet.go/rpc"
	"log"
	"math/big"
	"starknet/tool"
	"time"
)

type blockRefresh struct {
	blockNumber     chan uint64
	lastBlockNumber uint64
	ctx             context.Context
	cancel          context.CancelFunc
}

func InitBLock(lastBlockNumber uint64) {
	var br blockRefresh
	br.blockNumber = make(chan uint64, 5)
	br.lastBlockNumber = lastBlockNumber
	if br.lastBlockNumber == 0 {
		if blockNumber, err := getBlock(); err == nil {
			br.lastBlockNumber = blockNumber
		} else {
			log.Fatalln("getBlock ERR:%v", err)
		}
	}
	br.ctx, br.cancel = context.WithCancel(context.Background())
	go (&br).blockServer()
	(&br).hashServer()
}

func (b *blockRefresh) blockServer() {
	b.setBlock(b.lastBlockNumber)
	t := time.NewTicker(1 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-b.ctx.Done():
			return
		default:
			block, err := getBlock()
			if err != nil {
				b.cancel()
				continue
			}
			if block > b.lastBlockNumber {
				b.lastBlockNumber++
				b.setBlock(b.lastBlockNumber)
			}
			<-t.C
		}
	}
}

func (b *blockRefresh) setBlock(block uint64) {
	b.blockNumber <- block
}

func (b *blockRefresh) hashServer() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case blockNumber := <-b.blockNumber:
			blockData, err := getBlockWithTxs(blockNumber)
			if err != nil {
				fmt.Println("get blockWithTxs Err:", err)
				b.cancel()
				continue
			}
			fmt.Println("当前区块：", blockNumber)
			switch op := blockData.(type) {
			case *rpc.Block:
				fmt.Println("打包时间：", time.Unix(int64(op.Timestamp), 0).String())
				fmt.Println("交易笔数 ：", len(op.Transactions))
				for _, transaction := range op.Transactions {
					fmt.Println("交易hash：", transaction.Hash())
					transactionReceipt := getTransactionReceipt(transaction.Hash().String())
					fmt.Println("交易结果：", transactionReceipt.Result.ExecutionStatus)
					if fee, ok := new(big.Int).SetString(transactionReceipt.Result.ActualFee.Amount[2:], 16); ok {
						fmt.Println("手续费：", tool.ConvertWeiToUnit(fee, transactionReceipt.Result.ActualFee.Unit))
					}
					for _, v := range transactionReceipt.Result.Events {
						if v.Keys[0] == "0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9" {
							_, symbol, decimals := getTokenInfo(v.FromAddress)
							var _amount *big.Int
							if len(v.Data) == 2 {
								_amount = new(big.Int).SetInt64(0)
							} else {
								_amount, _ = new(big.Int).SetString(v.Data[2][2:], 16)
							}
							if decimals == 0 {
								if isNft721(v.FromAddress) {
									fmt.Println(fmt.Sprintf("FROM:%v  TO:%v   TokenId:%v %s", v.Data[0], v.Data[1], _amount.String(), symbol))
									continue
								} else if isNft1155(v.FromAddress) {
									fmt.Println(fmt.Sprintf("FROM:%v  TO:%v   TokenId:%v %s", v.Data[0], v.Data[1], _amount.String(), symbol))
								} else {
									fmt.Println(fmt.Sprintf("FROM:%v  TO:%v   未证实的721/1155 TokenId:%v %s", v.Data[0], v.Data[1], _amount.String(), symbol))
								}
							}
							divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
							amount := new(big.Float).Quo(new(big.Float).SetInt(_amount), new(big.Float).SetInt(divisor)).Text('f', decimals)
							fmt.Println(fmt.Sprintf("FROM:%v  TO:%v   AMOUNT:%v %s", v.Data[0], v.Data[1], amount, symbol))
						}
					}
					fmt.Println("====================新的交易================================")
				}
			case *rpc.PendingBlock:
			}
			b.lastBlockNumber = blockNumber
		}
	}
}
