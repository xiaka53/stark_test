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

func Transfer(account_addr, public_key, privateKey, contractAddr, _func string) {
	account_address, err := utils.HexToFelt(account_addr)
	if err != nil {
		panic(err.Error())
	}
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		panic(err.Error())
	}

	ks.Put(public_key, fakePrivKeyBI)

	maxfee, err := utils.HexToFelt("0x9184e72a000")
	if err != nil {
		panic(err.Error())
	}
	accnt, err := account.NewAccount(client, account_address, public_key, ks, 0)
	if err != nil {
		panic(err.Error())
	}
	nonce, rpcErr := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
	if rpcErr != nil {
		panic(rpcErr)
	}
	InvokeTx := rpc.InvokeTxnV1{
		MaxFee:        maxfee,
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: accnt.AccountAddress,
	}
	contractAddress, err := utils.HexToFelt(contractAddr)
	if err != nil {
		panic(err.Error())
	}
	amount, err := utils.HexToFelt("0x2")
	if err != nil {
		panic(err.Error())
	}
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(_func),
		Calldata:           []*felt.Felt{account_address, amount},
	}
	InvokeTx.Calldata, err = accnt.FmtCalldata([]rpc.FunctionCall{FnCall})
	if err != nil {
		panic(err.Error())
	}
	err = accnt.SignInvokeTransaction(context.Background(), &InvokeTx)
	if err != nil {
		panic(err.Error())
	}
	resp, rpcErr := accnt.AddInvokeTransaction(context.Background(), InvokeTx)
	if rpcErr != nil {
		panic(rpcErr)
	}
	fmt.Println("Transaction hash response : ", resp.TransactionHash)
}
