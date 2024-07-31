package main

import (
	"starknet/chain"
)

func main() {
	// 刷区块-解析数据
	// chain.InitBLock(0)
	// 地址部署
	chain.PrecomputeAddress("0x0253d0ce3d4c3ad9a6c4e86abf28b6193ef8799955f9bed8f7770464290d7d1f", "")
	// 执行交易
	// chain.Transfer("0x0253d0ce3d4c3ad9a6c4e86abf28b6193ef8799955f9bed8f7770464290d7d1f", "0x0253d0ce3d4c3ad9a6c4e86abf28b6193ef8799955f9bed8f7770464290d7d1f", "", "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d", "transfer")
}
