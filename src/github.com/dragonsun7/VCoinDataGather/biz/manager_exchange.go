package biz

import (
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"github.com/dragonsun7/VCoinDataGather/plugin"
	"github.com/dragonsun7/VCoinDataGather/model"
	"sync"
)

type ExchangeMgr struct {
}

var (
	exchangeMgrInstance *ExchangeMgr
	exchangeMgrOnce     sync.Once
)

func GetExchangeMgrInstance() (*ExchangeMgr) {
	exchangeMgrOnce.Do(func() {
		exchangeMgrInstance = new(ExchangeMgr)
	})

	return exchangeMgrInstance
}

// 从数据库中加载数据
func (em *ExchangeMgr) LoadData() ([]model.Exchange, error) {
	sql := `
SELECT
	uid, symbol, name_en, name_cn, website
FROM
	bs_exchange
WHERE
	active
`
	pg := postgres.GetInstance()
	dataSet, err := pg.Query(sql)
	if err != nil {
		return nil, err
	}

	var exchanges []model.Exchange
	for _, rec := range dataSet {
		var exchange model.Exchange
		uuid := rec["uid"].([]uint8)
		exchange.ID = string(uuid[:])
		exchange.Symbol = rec["symbol"].(string)
		exchange.NameEN = rec["name_en"].(string)
		exchange.NameCN = rec["name_cn"].(string)
		exchange.Website = rec["website"].(string)
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}

// 获取对应的插件
func (em *ExchangeMgr) GetPlugin(exchangeSymbol string) (plugin.Inf) {
	if "gateio" == exchangeSymbol {
		return new(plugin.GateioPlugin)
	}

	return nil
}
