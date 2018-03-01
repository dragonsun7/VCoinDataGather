package worker

import (
	"sync"
	"fmt"
	"time"
	"github.com/dragonsun7/VCoinDataGather/biz"
	"github.com/dragonsun7/VCoinDataGather/model"
	"github.com/dragonsun7/VCoinDataGather/lib"
	"github.com/dragonsun7/VCoinDataGather/config"
)

type TradeHistoryWorker struct {
	BaseWorker
	Exchange *model.Exchange
	Pair     *model.Pair
}

func (tw *TradeHistoryWorker) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	count := 1
	for true {
		if tw.isCancel {
			break
		}

		tradeHistoryMgr := biz.TradeHistoryMgr{Exchange: tw.Exchange, Pair: tw.Pair}

		tradeHistories, err := tradeHistoryMgr.GetData()
		if err != nil {
			lib.Logger().Println("\n重新请求，count：", count, tw.Pair.Symbol)
			fmt.Printf("获取交易历史(失败) %s, %s，第%d次 ...\n",
				tw.Exchange.Symbol, tw.Pair.Symbol, count)
			continue
		}

		err = tradeHistoryMgr.SaveData(tradeHistories)
		if err != nil {
			lib.Logger().Println("\n保存交易历史数据失败！", err)
			fmt.Printf("保存交易历史(失败) %s, %s，第%d次失败 ...\n",
				tw.Exchange.Symbol, tw.Pair.Symbol, count)
			break
		}

		fmt.Printf("获取交易历史(成功) %s, %s，第%d次，获取%d条 ...\n",
			tw.Exchange.Symbol, tw.Pair.Symbol, count, len(tradeHistories))


		cfg := config.GetInstance()
		seconds := time.Duration(cfg.Interval.TradeHistory)
		time.Sleep(time.Second * seconds)

		count++
	}
}
