package biz

import (
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"github.com/dragonsun7/VCoinDataGather/plugin"
	"github.com/dragonsun7/VCoinDataGather/model"
	"sync"
	"errors"
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

func (em *ExchangeMgr) GetExchange(symbol string) (model.Exchange, error) {
	var exchange model.Exchange

	sql := `SELECT uid, symbol, name_en, name_cn, website FROM bs_exchange WHERE active AND symbol = $1`
	pg := postgres.GetInstance()
	dataSet, err := pg.Query(sql, symbol)
	if err != nil {
		return exchange, err
	}
	if len(dataSet) != 1 {
		err = errors.New("交易所不存在、未激活或者表数据错误！")
		return exchange, err
	}

	rec := dataSet[0]
	uuid := rec["uid"].([]uint8)
	exchange.ID = string(uuid[:])
	exchange.Symbol = rec["symbol"].(string)
	exchange.NameEN = rec["name_en"].(string)
	exchange.NameCN = rec["name_cn"].(string)
	exchange.Website = rec["website"].(string)

	return exchange, nil
}

// 获取对应的插件
func (em *ExchangeMgr) GetPlugin(exchangeSymbol string) (plugin.Inf) {
	if "gateio" == exchangeSymbol {
		return new(plugin.GateioPlugin)
	}

	return nil
}
