package model

type TradeHistory struct {
	ID        string
	PairID    string
	TradeID   int64
	TimeStamp int64
	IsBuy     bool
	Price     float64 // 币种单价
	Quantity  float64 // 成交币种数量
}
