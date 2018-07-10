package core

import "core/types"

// 接收交易验证后加入交易池
// 需要排序根据时间
type TXPool struct {
	chain   *Blockchain
	pending map[string]*types.Transaction
}
