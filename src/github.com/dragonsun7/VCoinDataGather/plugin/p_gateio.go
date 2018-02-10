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
		return nil, err
	}

	resultStr, err := jsonparser.GetString(body, "result")
	if err != nil {
		return nil, err
	}

	result, err := strconv.ParseBool(resultStr)
	if err != nil {
		return nil, err
	}

	if !result {
		return nil, errors.New(resultStr)
	}

	var tradeHistories []model.TradeHistory
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var history model.JsonTradeHistory
		json.Unmarshal(value, &history)
		tradeHistory := history.ToTradeHistory()
		tradeHistories = append(tradeHistories, tradeHistory)
	}, "data")
	if err != nil {
		return nil, err
	}

	return tradeHistories, nil
}
