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
	"crypto/hmac"
	"crypto/sha512"
	"io/ioutil"
	"net/http"
)

type GateioPlugin struct {
	Plugin
}

// 获取交易对
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

// 获取交易历史
func (gateio *GateioPlugin) GetTradeHistory(pair string, lastID int64) ([]model.TradeHistory, error) {
	url := fmt.Sprintf("http://data.gate.io/api2/1/tradeHistory/%s/%d", pair, lastID)

	body, err := lib.HttpGet(url)
	if err != nil {
		lib.Logger().Println("\nHTTP请求失败：", err, url)
		return nil, err
	}

	err = gateio.parseResult(body)
	if err != nil {
		return nil, err
	}

	jsonStr := string(body)
	var tradeHistories []model.TradeHistory
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		tradeHistory, err1 := gateio.jsonToTradeHistory(value)
		if err1 != nil {
			return
		}

		tradeHistories = append(tradeHistories, *tradeHistory)
	}, "data")
	if err != nil {
		lib.Logger().Println("\n解析交易历史请求数据失败：", err, jsonStr)
		return nil, err
	}

	return tradeHistories, nil
}

// 获取用户历史交易记录(24小时)
func (gateio *GateioPlugin) GetUserTradeHistory(pair string, apiKey string, apiSecret string) ([]model.UserTradeHistory, error) {
	url := "https://api.gate.io/api2/1/private/tradeHistory"
	params := "currencyPair=" + pair

	body, err := gateio.post(url, params, apiKey, apiSecret)
	if err != nil {
		lib.Logger().Println("\nHTTP请求失败：", err, url)
		return nil, err
	}

	err = gateio.parseResult(body)
	if err != nil {
		return nil, err
	}

	jsonStr := string(body)
	var userTradeHistories []model.UserTradeHistory
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		userTradeHistory, err1 := gateio.jsonToUserTradeHistory(value)
		if err1 != nil {
			return
		}

		userTradeHistories = append(userTradeHistories, *userTradeHistory)
	}, "trades")
	if err != nil {
		lib.Logger().Println("\n解析用户交易历史请求数据失败：", err, jsonStr)
		return nil, err
	}

	return userTradeHistories, nil
}

/* -------------------- 私有函数 -------------------- */

func (gateio *GateioPlugin) createSign(params, secret string) (string) {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(params))
	sign := fmt.Sprintf("%x", mac.Sum(nil))

	return sign
}

func (gateio *GateioPlugin) post(url, params, key, secret string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", key)
	req.Header.Set("sign", gateio.createSign(params, secret))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// 解析result
func (gateio *GateioPlugin) parseResult(body []byte) (error) {
	jsonStr := string(body)

	resultStr, err := jsonparser.GetString(body, "result")
	if err != nil {
		lib.Logger().Println("\n从json解析result字段失败：", err, jsonStr)
		return err
	}

	result, err := strconv.ParseBool(resultStr)
	if err != nil {
		lib.Logger().Println("\n转换result字段值失败：", err, resultStr)
		return err
	}

	if !result {
		lib.Logger().Println("\n请求返回的结果值不为真：", err, resultStr)
		return errors.New(resultStr)
	}

	return nil
}

// 解析单条交易记录
func (gateio *GateioPlugin) jsonToTradeHistory(value []byte) (*model.TradeHistory, error) {
	var tradeHistory model.TradeHistory
	var err error
	jsonStr := string(value)

	// tradeID
	tradeHistory.TradeID, err = jsonparser.GetInt(value, "tradeID")
	if err != nil {
		lib.Logger().Println("\n解析tradeID失败：", err, jsonStr)
		return nil, err
	}

	// timestamp
	tradeHistory.TimeStamp, err = jsonparser.GetInt(value, "timestamp")
	if err != nil {
		lib.Logger().Println("\n解析timestamp失败：", err, jsonStr)
		return nil, err
	}

	// isBuy
	typeStr, err := jsonparser.GetString(value, "type")
	if err != nil {
		lib.Logger().Println("\n解析type失败：", err, jsonStr)
		return nil, err
	}
	tradeHistory.IsBuy = typeStr == "buy"

	// price
	tradeHistory.Price, err = jsonparser.GetFloat(value, "rate")
	if err != nil {
		lib.Logger().Println("\n解析rate失败：", err, jsonStr)
		return nil, err
	}

	// quantity
	tradeHistory.Quantity, err = jsonparser.GetFloat(value, "amount")
	if err != nil {
		tradeHistory.Quantity, err = lib.JsonStringValueToFloat64(value, "amount")
		if err != nil {
			lib.Logger().Println("\n解析amount失败：", err, jsonStr)
			return nil, err
		}
	}

	return &tradeHistory, nil
}

// 解析单条用户交易记录
func (gateio *GateioPlugin) jsonToUserTradeHistory(value []byte) (*model.UserTradeHistory, error) {
	var userTradeHistory model.UserTradeHistory
	var err error
	jsonStr := string(value)

	// tradeID
	userTradeHistory.TradeID, err = lib.JsonStringValueToInt64(value, "tradeID")
	if err != nil {
		lib.Logger().Println("\n解析tradeID失败：", err, jsonStr)
		return nil, err
	}

	// orderID
	userTradeHistory.OrderID, err = lib.JsonStringValueToInt64(value, "orderNumber")
	if err != nil {
		lib.Logger().Println("\n解析orderID失败：", err, jsonStr)
		return nil, err
	}

	// timestamp
	userTradeHistory.TimeStamp, err = lib.JsonStringValueToInt64(value, "time_unix")
	if err != nil {
		lib.Logger().Println("\n解析timestamp失败：", err, jsonStr)
		return nil, err
	}

	// isBuy
	typeStr, err := jsonparser.GetString(value, "type")
	if err != nil {
		lib.Logger().Println("\n解析type失败：", err, jsonStr)
		return nil, err
	}
	userTradeHistory.IsBuy = typeStr == "buy"

	// price
	userTradeHistory.Price, err = lib.JsonStringValueToFloat64(value, "rate")
	if err != nil {
		lib.Logger().Println("\n解析rate失败：", err, jsonStr)
		return nil, nil
	}

	// quantity (一般是浮点数，如果值是0则为字符串)
	userTradeHistory.Quantity, err = lib.JsonStringValueToFloat64(value, "amount")
	if err != nil {
		lib.Logger().Println("\n解析amount失败：", err, jsonStr)
		return nil, nil
	}

	// total
	userTradeHistory.Total, err = jsonparser.GetFloat(value, "total")
	if err != nil {
		lib.Logger().Println("\n解析total失败：", err, jsonStr)
		return nil, nil
	}

	return &userTradeHistory, nil
}
