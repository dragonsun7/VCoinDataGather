package biz

import (
	"github.com/dragonsun7/VCoinDataGather/model"
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
)

type TradeHistoryMgr struct {
	Exchange *model.Exchange
	Pair     *model.Pair
}

func (tm *TradeHistoryMgr) GetData() ([]model.TradeHistory, error) {
	exchangeMgr := GetExchangeMgrInstance()
	plugin := exchangeMgr.GetPlugin(tm.Exchange.Symbol)
	if plugin == nil {
		return nil, nil
	}

	return plugin.GetTradeHistory(tm.Pair.Symbol, tm.Pair.LastTradeHistoryID)
}

func (tm *TradeHistoryMgr) SaveData(tradeHistories []model.TradeHistory) (error) {
	pg := postgres.GetInstance()

	sql1 := `
INSERT INTO
  tr_history (pair_id, trade_id, isbuy, price, amount, ts)
VALUES
  ($1, $2, $3, $4, $5, to_timestamp($6))
ON
  CONFLICT(pair_id, trade_id)
DO
  NOTHING
`
	for _, tradeHistory := range tradeHistories {
		_, err := pg.Execute(sql1, tm.Pair.ID, tradeHistory.TradeID, tradeHistory.IsBuy,
			tradeHistory.Price, tradeHistory.Amount, tradeHistory.TimeStamp)
		if err != nil {
			return err
		}
	}

	// 更新last_th_id
	sql2 := `
UPDATE
	bs_pair
SET
	last_th_id = $1
WHERE
	uid = $2
`
	lastIndex := len(tradeHistories) - 1
	if lastIndex >= 0 {
		tm.Pair.LastTradeHistoryID = tradeHistories[lastIndex].TradeID
		_, err := pg.Execute(sql2, tm.Pair.LastTradeHistoryID, tm.Pair.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
