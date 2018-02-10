// 获取交易对信息

package main

import (
	"fmt"
	"github.com/dragonsun7/VCoinDataGather/biz"
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
)

func main() {
	defer postgres.GetInstance().CloseDB()

	fmt.Println("开始获取可用的交易所...")
	exchangeMgr := biz.GetExchangeMgrInstance()
	exchanges, err := exchangeMgr.LoadData()
	if err != nil {
		panic(err)
	}

	for _, exchange := range exchanges {
		pairMgr := biz.PairMgr{Exchange:&exchange}
		fmt.Printf("开始获取%s的交易对...\n", exchange.Symbol)

		pairs, err := pairMgr.GetData()
		if err != nil {
			panic(err)
		}

		err = pairMgr.SaveData(pairs)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("执行完毕！")
}
