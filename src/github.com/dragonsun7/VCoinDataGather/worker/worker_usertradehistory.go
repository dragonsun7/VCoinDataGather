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
	prevTime := time.Now()
	for true {
		if worker.isCancel {
			break
		}

		// 每天00:00:00更新一次
		prevHour := prevTime.Hour()
		nowHour := time.Now().Hour()
		prevTime = time.Now()
		if (nowHour == prevHour) || (nowHour != 0) {
			continue
		}

		userTradeHistoryMgr := biz.UserTradeHistoryMgr{
			Exchange: worker.Exchange, User: worker.User, Pair: worker.Pair, API: worker.API}

		userTradeHistories, err := userTradeHistoryMgr.GetData()
		if err != nil {
			lib.Logger().Println("\n重新请求，count：", count)
			fmt.Printf("获取用户交易历史(失败) %s, %s, %s，第%d次 ...\n",
				worker.User.Username, worker.Exchange.Symbol, worker.Pair.Symbol, count)
			continue
		}

		err = userTradeHistoryMgr.SaveData(userTradeHistories)
		if err != nil {
			lib.Logger().Println("\n保存用户交易历史数据失败！", err)
			fmt.Printf("保存用户交易历史(失败) %s, %s, %s，第%d次 ...\n",
				worker.User.Username, worker.Exchange.Symbol, worker.Pair.Symbol, count)
			break
		}

		fmt.Printf("获取用户交易历史(成功) %s, %s, %s，第%d次, 获取%d条 ...\n",
			worker.User.Username, worker.Exchange.Symbol, worker.Pair.Symbol, count, len(userTradeHistories))

		cfg := config.GetInstance()
		seconds := time.Duration(cfg.Interval.UserTradeHistory)
		time.Sleep(time.Second * seconds)

		count++
	}
}
