package plugin

import (
	"github.com/dragonsun7/VCoinDataGather/model"
)

type Inf interface {
	// 获取交易对
	GetPairs() ([]string, error)

	// 获取历史交易记录
	GetTradeHistory(pair string, lastID int64) ([]model.TradeHistory, error)

	// 获取用户历史交易记录
	GetUserTradeHistory(pair string, apiKey string, apiSecret string) ([]model.UserTradeHistory, error)
}

type Plugin struct {
}
