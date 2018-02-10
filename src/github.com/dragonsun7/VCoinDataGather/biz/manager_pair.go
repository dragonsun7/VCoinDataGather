package biz

import (
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"strings"
	"errors"
	"github.com/dragonsun7/VCoinDataGather/model"
)

type PairMgr struct {
	Exchange *model.Exchange
}

// 通过平台接口获取数据
func (pm *PairMgr) GetData() ([]model.Pair, error) {
	exchangeMgr := GetExchangeMgrInstance()
	plugin := exchangeMgr.GetPlugin(pm.Exchange.Symbol)
	if plugin == nil {
		return nil, nil
	}

	pairsArr, err := plugin.GetPairs()
	if err != nil {
		return nil, err
	}

	var pairs []model.Pair
	for _, symbol := range pairsArr {
		var pair model.Pair
		currs := strings.Split(symbol, "_")
		if len(currs) != 2 {
			return nil, errors.New("交易对解析错误！")
		}
		pair.Symbol = strings.ToUpper(symbol)
		pair.CurrA = strings.ToUpper(currs[0])
		pair.CurrB = strings.ToUpper(currs[1])
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

// 从数据库中加载数据
func (pm *PairMgr) LoadData() ([]model.Pair, error) {
	sql := `
SELECT
	uid, symbol, curr_a, curr_b, last_th_id
FROM
	bs_pair
WHERE
	exchange_id = $1
`

	pg := postgres.GetInstance()
	dataSet, err := pg.Query(sql, pm.Exchange.ID)
	if err != nil {
		return nil, err
	}

	var pairs []model.Pair
	for _, rec := range dataSet {
		var pair model.Pair
		uuid := rec["uid"].([]uint8)
		pair.ID = string(uuid[:])
		pair.Symbol = rec["symbol"].(string)
		pair.CurrA = rec["curr_a"].(string)
		pair.CurrB = rec["curr_b"].(string)
		pair.LastTradeHistoryID = rec["last_th_id"].(int64)
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

// 保存数据到数据库
func (pm *PairMgr) SaveData(pairs []model.Pair) (error) {
	sql := `
INSERT INTO
	bs_pair (exchange_id, symbol, curr_a, curr_b)
VALUES
	($1, $2, $3, $4)
ON
 	conflict(exchange_id, symbol)
DO 
	NOTHING
`
	pg := postgres.GetInstance()

	for _, pair := range pairs {
		_, err := pg.Execute(sql, pm.Exchange.ID, pair.Symbol, pair.CurrA, pair.CurrB)
		if err != nil {
			return err
		}
	}

	return nil
}
