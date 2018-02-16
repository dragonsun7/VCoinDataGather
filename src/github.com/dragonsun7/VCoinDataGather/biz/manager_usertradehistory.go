package biz

import (
	"github.com/dragonsun7/VCoinDataGather/model"
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
)

type UserTradeHistoryMgr struct {
	Exchange *model.Exchange
	User     *model.User
	Pair     *model.Pair
	API      *model.API
}

func (mgr *UserTradeHistoryMgr) GetData() ([]model.UserTradeHistory, error) {
	exchangeMgr := GetExchangeMgrInstance()
	plugin := exchangeMgr.GetPlugin(mgr.Exchange.Symbol)
	if plugin == nil {
		return nil, nil
	}

	return plugin.GetUserTradeHistory(mgr.Pair.Symbol, mgr.API.Key, mgr.API.Secret)
}

func (mgr *UserTradeHistoryMgr) SaveData(userTradeHistories []model.UserTradeHistory) (error) {
	pg := postgres.GetInstance()

	sql := `
INSERT INTO
	tr_usr_history (pair_id, user_id, trade_id, order_id, isbuy, price, quantity, total, ts)
VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, to_timestamp($9))
ON
	CONFLICT(pair_id, user_id, trade_id, order_id)
DO
	NOTHING
`
	for _, userTradeHistory := range userTradeHistories {
		_, err := pg.Execute(sql,
			mgr.Pair.ID, mgr.User.ID, userTradeHistory.TradeID, userTradeHistory.OrderID, userTradeHistory.IsBuy,
			userTradeHistory.Price, userTradeHistory.Quantity, userTradeHistory.Total, userTradeHistory.TimeStamp)
		if err != nil {
			return err
		}
	}

	return nil
}
