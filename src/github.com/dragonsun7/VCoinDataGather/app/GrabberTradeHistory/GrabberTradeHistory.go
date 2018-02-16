package main

import (
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"github.com/dragonsun7/VCoinDataGather/biz"
	"github.com/dragonsun7/VCoinDataGather/worker"
	"github.com/dragonsun7/VCoinDataGather/lib"
	"fmt"
	"sync"
)

const (
	logfile = "GrabberTradeHistory.log"
)

func main() {
	defer postgres.GetInstance().CloseDB()

	// 初始化日志
	err := lib.LoggerInit(logfile)
	if err != nil {
		panic(err)
	}
	defer lib.LoggerClose()

	// 业务处理
	fmt.Println("开始获取可用的交易所...")
	exchangeMgr := biz.GetExchangeMgrInstance()
	exchanges, err := exchangeMgr.LoadData()
	if err != nil {
		panic(err)
	}

	waitGroup := &sync.WaitGroup{}
	for _, exchange := range exchanges {
		theExchange := exchange
		pairMgr := biz.PairMgr{Exchange: &theExchange}
		pairs, err := pairMgr.LoadData()
		if err != nil {
			panic(err)
		}

		fmt.Printf("开始获取%s的交易对...\n", exchange.Symbol)
		for _, pair := range pairs {
			thePair := pair
			tradeHistoryWorker := worker.TradeHistoryWorker{Exchange: &theExchange, Pair: &thePair}
			waitGroup.Add(1)
			go tradeHistoryWorker.Start(waitGroup)
		}
	}
	waitGroup.Wait()

	fmt.Println("完成！")
}
