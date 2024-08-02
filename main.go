package main

import (
	"fmt"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"starknet/chain"
	"time"
)

var (
	EthContentAddress, DeployAddressFeeAccountPri *felt.Felt
)

func init() {
	var err error
	EthContentAddress, err = utils.HexToFelt("0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7")
	if err != nil {
		panic(fmt.Sprintf("ethcontentAddress is err"))
	}
	DeployAddressFeeAccountPri, err = utils.HexToFelt("0x00ea425b8b1fa8c674b8291611044aa7d9c6b3e79ca70c8ddbc22a8b51b81f0b")
	if err != nil {
		panic(fmt.Sprintf("deployAddressFeeAccountPri is err"))
	}
	// DeployAddressFeeAccountAddr, err = utils.HexToFelt("0x006264b39ee4eb5ce1e7d3a56dcc2b25fa4facb67056f3a81786e49c726c4018")
	// if err != nil {
	// 	panic(fmt.Sprintf("ethcontentAddress is err"))
	// }
}

func main() {
	// 刷区块-解析数据
	// chain.InitBLock(0)

	// 地址生成
	pubAddr, priAddr := chain.CreateAddress()
	fmt.Println(fmt.Sprintf("pubAddr:	%s", pubAddr))
	fmt.Println(fmt.Sprintf("priAddr:	%s", priAddr))

	// 执行交易
	chain.Transfer(DeployAddressFeeAccountPri, pubAddr, EthContentAddress, new(felt.Felt).SetUint64(500000000000000), "transfer")
	fmt.Println("转出成功，等待链上校验后进行账号部署")
	time.Sleep(5 * time.Minute)
	// 地址部署
	chain.DeployAddress(priAddr)
}
