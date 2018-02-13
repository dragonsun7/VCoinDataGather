package plugin

import (
	"github.com/dragonsun7/VCoinDataGather/lib"
	"github.com/dragonsun7/VCoinDataGather/model"
	"github.com/buger/jsonparser"
	"encoding/json"
	"strings"
	"fmt"
	"strconv"
	"errors"
)

type GateioPlugin struct {
	Plugin
}

func (gateio *GateioPlugin) GetPairs() ([]string, error) {
	url := "http://data.gate.io/api2/1/pairs"

	body, err := lib.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var pairs []string
	err = json.Unmarshal(body, &pairs)
	if err != nil {
		return nil, err
	}

	// 暂时只要USDT的交易对
	var ret []string
	for _, pair := range pairs {
		if strings.HasSuffix(strings.ToUpper(pair), "USDT") {
			ret = append(ret, pair)
		}
	}

	return ret, nil
}

func (gateio *GateioPlugin) GetTradeHistory(pair string, lastID int64) ([]model.TradeHistory, error) {
	url := fmt.Sprintf("http://data.gate.io/api2/1/tradeHistory/%s/%d", pair, lastID)

	body, err := lib.HttpGet(url)
	if err != nil {
		lib.Logger().Println("\nHTTP请求失败：", err, url)
		return nil, err
	}

	resultStr, err := jsonparser.GetString(body, "result")
	if err != nil {
		lib.Logger().Println("\n从json解析result字段失败：", err, string(body[:]))
		return nil, err
	}

	result, err := strconv.ParseBool(resultStr)
	if err != nil {
		lib.Logger().Println("\n转换result字段值失败：", err, resultStr)
		return nil, err
	}

	if !result {
		lib.Logger().Println("\n请求返回的结果值不为真：", err, resultStr)
		return nil, errors.New(resultStr)
	}

	var tradeHistories []model.TradeHistory
	var err1 error
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var tradeHistory model.TradeHistory
		json := string(value)

		tradeHistory.TradeID, err1 = jsonparser.GetInt(value, "tradeID")
		if err1 != nil {
			lib.Logger().Println("\n解析tradeID失败：", err1, json)
			return
		}

		tradeHistory.TimeStamp, err1 = jsonparser.GetInt(value, "timestamp")
		if err1 != nil {
			lib.Logger().Println("\n解析timestamp失败：", err1, json)
			return
		}

		var typeStr string
		typeStr, err1 = jsonparser.GetString(value, "type")
		if err1 != nil {
			lib.Logger().Println("\n解析type失败：", err1, json)
			return
		}
		tradeHistory.IsBuy = typeStr == "buy"

		tradeHistory.Price, err1 = jsonparser.GetFloat(value, "rate")
		if err1 != nil {
			lib.Logger().Println("\n解析rate失败：", err1, json)
			return
		}

		tradeHistory.Quantity, err1 = jsonparser.GetFloat(value, "amount")
		if err1 != nil {
			var amountStr string

			amountStr, err1 = jsonparser.GetString(value, "amount")
			if err1 != nil {
				lib.Logger().Println("\n解析amount失败：", err1, json)
				return
			}

			tradeHistory.Quantity, err1 = strconv.ParseFloat(amountStr, 64)
			if err1 != nil {
				lib.Logger().Println("\n转换为浮点数失败：", err1, amountStr)
				return
			}
		}

		tradeHistories = append(tradeHistories, tradeHistory)
	}, "data")
	if err != nil {
		lib.Logger().Println("\n解析交易历史请求数据失败：", err, string(body[:]))
		return nil, err
	}

	return tradeHistories, nil
}
