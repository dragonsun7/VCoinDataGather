package main

import (
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"github.com/dragonsun7/VCoinDataGather/lib"
	"fmt"
	"github.com/dragonsun7/VCoinDataGather/biz"
	worker2 "github.com/dragonsun7/VCoinDataGather/worker"
	"sync"
)

const (
	logfile      = "GrabberUserTradeHistory.log"
	username     = "dragonsun7"
	exchangeName = "gateio"
)

func main() {
	defer postgres.GetInstance().CloseDB()

	/* 初始化日志 */
	err := lib.LoggerInit(logfile)
	if err != nil {
		panic(err)
	}
	defer lib.LoggerClose()

	/* 业务处理 */
	user, err := biz.GetUserMgrInstance().GetUser(username)
	if err != nil {
		fmt.Println("获取用户信息失败！")
		panic(err)
	}

	exchange, err := biz.GetExchangeMgrInstance().GetExchange(exchangeName)
	if err != nil {
		fmt.Println("获取交易所信息失败！")
		panic(err)
	}

	api, err := biz.GetAPIMgrInstance().GetAPI(user, exchange)
	if err != nil {
		fmt.Println("获取API信息失败！")
		panic(err)
	}

	pairMgr := biz.PairMgr{Exchange: &exchange}
	pairs, err := pairMgr.LoadData()
	if err != nil {
		fmt.Println("获取交易对失败！")
		panic(err)
	}

	waitGroup := &sync.WaitGroup{}
	for _, pair := range pairs {
		thePair := pair
		worker := worker2.UserTradeHistoryWorker{Exchange: &exchange, User: &user, Pair: &thePair, API: &api}
		waitGroup.Add(1)
		go worker.Start(waitGroup)
	}
	waitGroup.Wait()
}
