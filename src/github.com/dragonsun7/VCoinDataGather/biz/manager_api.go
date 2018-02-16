package biz

import (
	"sync"
	"github.com/dragonsun7/VCoinDataGather/model"
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"errors"
)

type APIMgr struct {
}

var (
	apiMgrInstance *APIMgr
	apiMgrOnce     sync.Once
)

func GetAPIMgrInstance() (*APIMgr) {
	apiMgrOnce.Do(func() {
		apiMgrInstance = new(APIMgr)
	})

	return apiMgrInstance
}

// 获取 API Key 和 API Secret
func (am *APIMgr) GetAPI(user model.User, exchange model.Exchange) (model.API, error) {
	var api model.API

	sql := `
SELECT
	a.uid,
	a.api_key,
	a.api_secret
FROM
	st_api AS a,
	bs_user AS u
WHERE
	a.user_id = u.uid
	AND u.active
	AND a.user_id = $1
	AND a.exchange_id = $2
`
	pg := postgres.GetInstance()
	dataSet, err := pg.Query(sql, user.ID, exchange.ID)
	if err != nil {
		return api, err
	}
	if len(dataSet) != 1 {
		err = errors.New("没有找到对应的API信息，或者表数据不正确！")
		return api, err
	}

	rec := dataSet[0]
	uuid := rec["uid"].([]uint8)
	api.ID = string(uuid[:])
	api.UserID = user.ID
	api.ExchangeID = exchange.ID
	api.Key = rec["api_key"].(string)
	api.Secret = rec["api_secret"].(string)

	return api, nil
}
