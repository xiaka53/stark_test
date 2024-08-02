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

// func Transfer(account_addr, public_key, privateKey, contractAddr, _func string) {
func Transfer(privateKey, toAddr, contractAddr, amount *felt.Felt, _func string) {
	public_key := exportPubKey(privateKey)
	account_addr, err := exportPubAddr(public_key, privateKey)

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey.Text(16), 16)
	if !ok {
		panic("privatekey is err")
	}

	ks.Put(public_key.String(), fakePrivKeyBI)
	accnt, err := account.NewAccount(client, account_addr, public_key.String(), ks, 2)
	if err != nil {
		panic(err.Error())
	}
	nonce, rpcErr := accnt.Nonce(context.Background(), rpc.BlockID{Tag: "latest"}, accnt.AccountAddress)
	if rpcErr != nil {
		panic(rpcErr)
	}
	InvokeTx := rpc.InvokeTxnV1{
		MaxFee:        new(felt.Felt).SetUint64(500000000000000),
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: accnt.AccountAddress,
	}
	amountDown := new(felt.Felt).SetUint64(0)
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractAddr,
		EntryPointSelector: utils.GetSelectorFromNameFelt(_func),
		Calldata:           []*felt.Felt{toAddr, amount, amountDown},
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
