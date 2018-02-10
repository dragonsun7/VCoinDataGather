/*
	轮询抓取各平台行情数据，并保存到数据库中
	多线程抓取不同平台的交易对(每个整点抓取一次)
	多线程抓取不同平台不同交易对的数据(每10秒抓取一次)
 */

package main

import (
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"github.com/dragonsun7/VCoinDataGather/biz"
	"github.com/dragonsun7/VCoinDataGather/worker"
	"fmt"
	"sync"
)

func main() {
	defer postgres.GetInstance().CloseDB()

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
