package model

// 交易所
type Exchange struct {
	ID      string
	Symbol  string // 符号(32, 不可为空)
	NameEN  string // 英文名(32, 可为空)
	NameCN  string // 中文名(32, 可为空)
	Website string // 网站地址(128, 可为空)
}
