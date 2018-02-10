package model

type TradeHistory struct {
	ID        string
	PairID    string
	TradeID   int64
	TimeStamp int64
	IsBuy     bool
	Price     float64
	Amount    float64
}

type JsonTradeHistory struct {
	TradeID   int64   `json:"tradeID"`   // tradeID
	TimeStamp int64   `json:"timestamp"` // timestamp
	Type      string  `json:"type"`      // 交易类型, buy买 sell卖
	Price     float64 `json:"rate"`      // 币种单价
	Amount    float64 `json:"amount"`    // 成交币种数量
	Total     float64 `json:"total"`     // 订单总额
	Date      string  `json:"date"`      // 订单时间
}

func (jth *JsonTradeHistory) ToTradeHistory() (TradeHistory) {
	var tradeHistory TradeHistory
	tradeHistory.TradeID = jth.TradeID
	tradeHistory.TimeStamp = jth.TimeStamp
	tradeHistory.IsBuy = jth.Type == "buy"
	tradeHistory.Price = jth.Price
	tradeHistory.Amount = jth.Amount

	return tradeHistory
}
