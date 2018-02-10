package plugin

import (
	"github.com/dragonsun7/VCoinDataGather/model"
)

type Inf interface {
	// 获取交易对
	GetPairs() ([]string, error)

	// 获取历史交易记录
	GetTradeHistory(pair string, lastID int64) ([]model.TradeHistory, error)
}

type Plugin struct {
}
