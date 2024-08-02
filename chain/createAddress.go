package chain

import (
	"context"
	"errors"
	"fmt"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"math/big"
)

var (
	predeployedClassHash = "0x00e2eb8f5672af4e6a4e8a8f1b44989685e668489b0a25437733756c5a34a1d6"
)

func DeployAddress(privateKey *felt.Felt) {
	pub := exportPubKey(privateKey)

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey.Text(16), 16)
	if !ok {
		panic("privatekey is err")
	}

	ks.Put(pub.String(), fakePrivKeyBI)

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
		MaxFee:              new(felt.Felt).SetUint64(500000000000000),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		ClassHash:           classHash,
		ContractAddressSalt: pub,
		ConstructorCalldata: []*felt.Felt{pub},
	}
	precomputeAddress, err := exportPubAddr(pub, privateKey)
	if err != nil {
		panic(err.Error())
	}
	err = accnt.SignDeployAccountTransaction(context.Background(), &tx, precomputeAddress)
	if err != nil {
		panic(err)
	}
	resp, rpcErr := accnt.AddDeployAccountTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{DeployAccountTxn: tx})
	if rpcErr != nil {
		panic(fmt.Sprint("Error returned from AddDeployAccountTransaction: ", rpcErr))
	}
	fmt.Println("AddDeployAccountTransaction response:", resp)
}

func CreateAddress() (*felt.Felt, *felt.Felt) {
	_, pub, pri := account.GetRandomKeys()
	addr, err := exportPubAddr(pub, pri)
	if err != nil {
		panic(fmt.Sprintf("create address err:%v", err))
	}
	return addr, pri
}

// 私钥导出公钥
func exportPubKey(pri *felt.Felt) *felt.Felt {
	privatekey, ok := new(big.Int).SetString(pri.Text(16), 16)
	if !ok {
		panic("privatekey is err")
	}
	pubx, puby, err := curve.Curve.PrivateToPoint(privatekey)
	if err != nil {
		return nil
	}
	if !curve.Curve.IsOnCurve(pubx, puby) {
		return nil
	}
	return new(felt.Felt).SetBytes(pubx.Bytes())
}

// 验证私钥公钥是否匹配
func verifyPriPub(pub, pri *felt.Felt) bool {
	getPub := exportPubKey(pri)
	if getPub == nil {
		return false
	}

	return getPub.String() == pub.String()
}

// 私钥导出抽象地址
func exportPubAddr(pub, pri *felt.Felt) (*felt.Felt, error) {
	if !verifyPriPub(pub, pri) {
		return nil, errors.New("private key public key does not match")
	}
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(pri.Text(16), 16)
	if !ok {
		return nil, errors.New("privekey is err")
	}
	ks.Put(pub.String(), fakePrivKeyBI)
	accnt, err := account.NewAccount(client, pub, pub.String(), ks, 2)
	if err != nil {
		return nil, err
	}
	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		return nil, err
	}
	return accnt.PrecomputeAddress(&felt.Zero, pub, classHash, []*felt.Felt{pub})
}
