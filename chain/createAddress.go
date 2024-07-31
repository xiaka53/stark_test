package chain

import (
	"context"
	"fmt"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"math/big"
)

var (
	predeployedClassHash = "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f"
)

func PrecomputeAddress(account_addr, privateKey string) {
	pub, err := utils.HexToFelt(account_addr)
	if err != nil {
		panic(err.Error())
	}
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic(err.Error())
	}

	ks.Put(privateKey, fakePrivKeyBI)

	accnt, err := account.NewAccount(client, pub, pub.String(), ks, 2)
	if err != nil {
		panic(err)
	}
	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		panic(err)
	}
	tx := rpc.DeployAccountTxn{
		Nonce:               &felt.Zero,
		MaxFee:              new(felt.Felt).SetUint64(7268996239700),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		ClassHash:           classHash,
		ContractAddressSalt: pub,
		ConstructorCalldata: []*felt.Felt{pub},
	}
	precomputedAddress, err := accnt.PrecomputeAddress(&felt.Zero, pub, classHash, tx.ConstructorCalldata)
	fmt.Println(precomputedAddress.String())
	if err != nil {
		panic(err)
	}
	err = accnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress)
	if err != nil {
		panic(err)
	}
	resp, rpcErr := accnt.AddDeployAccountTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{DeployAccountTxn: tx})
	if rpcErr != nil {
		panic(fmt.Sprint("Error returned from AddDeployAccountTransaction: ", rpcErr))
	}
	fmt.Println("AddDeployAccountTransaction response:", resp)
}

func CreateAddress() (string, string) {
	_, pub, pri := account.GetRandomKeys()
	pubAddr := pub.String()
	priAddr := pri.String()
	fmt.Println(pubAddr)
	fmt.Println(priAddr)

	fmt.Println(len(pubAddr), len(priAddr))
	return priAddr, priAddr
}
