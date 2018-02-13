package worker

import (
	"sync"
	"fmt"
	"time"
	"github.com/dragonsun7/VCoinDataGather/biz"
	"github.com/dragonsun7/VCoinDataGather/model"
	"github.com/dragonsun7/VCoinDataGather/lib"
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

		fmt.Printf("获取交易历史 %s, %s，第%d次 ...\n", tw.Exchange.Symbol, tw.Pair.Symbol, count)
		tradeHistoryMgr := biz.TradeHistoryMgr{Exchange: tw.Exchange, Pair: tw.Pair}
		tradeHistories, err := tradeHistoryMgr.GetData()
		if err != nil {
			lib.Logger().Println("重新请求，count：", count)
			continue
		}

		err = tradeHistoryMgr.SaveData(tradeHistories)
		if err != nil {
			lib.Logger().Println("\n保存交易历史数据失败！", err)
			break
		}

		time.Sleep(time.Second * 10)
		count++
	}
}
