package chain

import (
	"encoding/json"
	"math/big"
	"starknet/request"
)

// 获取最新区块高度
func getBlockNumber() uint64 {
	data, err := request.Post(getQuery(BLOCKNUMBER, ""))
	if err != nil {
		return 0
	}
	var _blockNumber blockkNumber
	if err = json.Unmarshal(data, &_blockNumber); err != nil {
		return 0
	}
	return _blockNumber.Result
}

// 获取区块交易
func getBlockByNumber(blockNumber *big.Int) ([]Transaction, int) {
	data, err := request.Post(getQuery(BLOCKTXS, []map[string]any{{"block_number": blockNumber}}))
	if err != nil {
		return nil, 0
	}
	var _block Block
	if err = json.Unmarshal(data, &_block); err != nil {
		return nil, 0
	}
	return _block.Result.Transactions, _block.Result.Timestamp
}

// 获取hash交易
func getTransactionReceipt(hash string) TransactionReceiptData {
	var _transactionReceipt TransactionReceiptData
	data, err := request.Post(getQuery(TransactionReceipt, []string{hash}))
	if err != nil {
		return _transactionReceipt
	}
	if err = json.Unmarshal(data, &_transactionReceipt); err != nil {
		return _transactionReceipt
	}
	return _transactionReceipt
}

func getcall(hash string) TransactionReceiptData {
	var _transactionReceipt TransactionReceiptData
	data, err := request.Post(getQuery(CALL, []string{hash}))
	if err != nil {
		return _transactionReceipt
	}
	if err = json.Unmarshal(data, &_transactionReceipt); err != nil {
		return _transactionReceipt
	}
	return _transactionReceipt
}
