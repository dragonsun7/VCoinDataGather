package model

type UserTradeHistory struct {
	ID        string
	UserID    string
	PairID    string
	TradeID   int64
	OrderID   int64
	TimeStamp int64
	IsBuy     bool
	Price     float64 // 币种单价
	Quantity  float64 // 成交币种数量
	Total     float64 // 总价
}
