package worker

import (
	"sync"
	"fmt"
	"time"
	"github.com/dragonsun7/VCoinDataGather/biz"
	"github.com/dragonsun7/VCoinDataGather/model"
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
			fmt.Println(err)
			break
		}

		err = tradeHistoryMgr.SaveData(tradeHistories)
		if err != nil {
			fmt.Print(err)
			break
		}

		time.Sleep(time.Second * 10)
		count++
	}
}
