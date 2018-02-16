package worker

import (
	"github.com/dragonsun7/VCoinDataGather/model"
	"sync"
	"fmt"
	"github.com/dragonsun7/VCoinDataGather/biz"
	"github.com/dragonsun7/VCoinDataGather/lib"
	"time"
	"github.com/dragonsun7/VCoinDataGather/config"
)

type UserTradeHistoryWorker struct {
	BaseWorker
	Exchange *model.Exchange
	User     *model.User
	Pair     *model.Pair
	API      *model.API
}

func (worker *UserTradeHistoryWorker) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	count := 1
	for true {
		if worker.isCancel {
			break
		}

		fmt.Printf("获取用户交易历史 %s, %s, %s，第%d次 ...\n",
			worker.User.Username, worker.Exchange.Symbol, worker.Pair.Symbol, count)

		userTradeHistoryMgr := biz.UserTradeHistoryMgr{
			Exchange: worker.Exchange, User: worker.User, Pair: worker.Pair, API: worker.API}

		userTradeHistories, err := userTradeHistoryMgr.GetData()
		if err != nil {
			lib.Logger().Println("\n重新请求，count：", count)
			continue
		}

		err = userTradeHistoryMgr.SaveData(userTradeHistories)
		if err != nil {
			lib.Logger().Println("\n保存用户交易历史数据失败！", err)
			break
		}

		cfg := config.GetInstance()
		seconds := time.Duration(cfg.Interval.UserTradeHistory)
		time.Sleep(time.Second * seconds)

		count++
	}
}
