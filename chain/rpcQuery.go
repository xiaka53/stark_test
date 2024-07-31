package chain

import "encoding/json"

const (
	BLOCKNUMBER        = "starknet_blockNumber"
	BLOCKTXS           = "starknet_getBlockWithTxs"
	TransactionReceipt = "starknet_getTransactionReceipt"
	CALL               = "starknet_call"
)

type blockkNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  uint64 `json:"result"`
	Id      int    `json:"id"`
}

type Block struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Status           string `json:"status"`
		BlockHash        string `json:"block_hash"`
		ParentHash       string `json:"parent_hash"`
		BlockNumber      int    `json:"block_number"`
		NewRoot          string `json:"new_root"`
		Timestamp        int    `json:"timestamp"`
		SequencerAddress string `json:"sequencer_address"`
		L1GasPrice       struct {
			PriceInFri string `json:"price_in_fri"`
			PriceInWei string `json:"price_in_wei"`
		} `json:"l1_gas_price"`
		StarknetVersion string        `json:"starknet_version"`
		Transactions    []Transaction `json:"transactions"`
	} `json:"result"`
	Id int `json:"id"`
}

type Transaction struct {
	TransactionHash string   `json:"transaction_hash"`
	Type            string   `json:"type"`
	Version         string   `json:"version"`
	Nonce           string   `json:"nonce"`
	MaxFee          string   `json:"max_fee,omitempty"`
	SenderAddress   string   `json:"sender_address,omitempty"`
	Signature       []string `json:"signature"`
	Calldata        []string `json:"calldata,omitempty"`
	ResourceBounds  struct {
		L1Gas struct {
			MaxAmount       string `json:"max_amount"`
			MaxPricePerUnit string `json:"max_price_per_unit"`
		} `json:"l1_gas"`
		L2Gas struct {
			MaxAmount       string `json:"max_amount"`
			MaxPricePerUnit string `json:"max_price_per_unit"`
		} `json:"l2_gas"`
	} `json:"resource_bounds,omitempty"`
	Tip                       string        `json:"tip,omitempty"`
	PaymasterData             []interface{} `json:"paymaster_data,omitempty"`
	AccountDeploymentData     []interface{} `json:"account_deployment_data,omitempty"`
	NonceDataAvailabilityMode string        `json:"nonce_data_availability_mode,omitempty"`
	FeeDataAvailabilityMode   string        `json:"fee_data_availability_mode,omitempty"`
	ContractAddressSalt       string        `json:"contract_address_salt,omitempty"`
	ClassHash                 string        `json:"class_hash,omitempty"`
	ConstructorCalldata       []string      `json:"constructor_calldata,omitempty"`
}

type TransactionReceiptData struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Type            string `json:"type"`
		TransactionHash string `json:"transaction_hash"`
		ActualFee       struct {
			Amount string `json:"amount"`
			Unit   string `json:"unit"`
		} `json:"actual_fee"`
		ExecutionStatus string        `json:"execution_status"`
		FinalityStatus  string        `json:"finality_status"`
		BlockHash       string        `json:"block_hash"`
		BlockNumber     int           `json:"block_number"`
		MessagesSent    []interface{} `json:"messages_sent"`
		Events          []struct {
			FromAddress string   `json:"from_address"`
			Keys        []string `json:"keys"`
			Data        []string `json:"data"`
		} `json:"events"`
		ExecutionResources struct {
			Steps                         int `json:"steps"`
			MemoryHoles                   int `json:"memory_holes"`
			PedersenBuiltinApplications   int `json:"pedersen_builtin_applications"`
			RangeCheckBuiltinApplications int `json:"range_check_builtin_applications"`
		} `json:"execution_resources"`
	} `json:"result"`
	Id int `json:"id"`
}

type rpcQuery struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      int    `json:"id"`
	Params  any    `json:"params"`
}

func getQuery(methon string, params any) []byte {
	query := rpcQuery{
		Jsonrpc: "2.0",
		Method:  methon,
		Id:      0,
	}
	if params != "" {
		query.Params = params
	}
	queryByte, _ := json.Marshal(query)
	return queryByte
}
