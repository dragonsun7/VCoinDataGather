package biz

import (
	"sync"
	"github.com/dragonsun7/VCoinDataGather/model"
	"github.com/dragonsun7/VCoinDataGather/db/postgres"
	"errors"
)

type UserMgr struct {
}

var (
	userMgrInstance *UserMgr
	userMgrOnce     sync.Once
)

func GetUserMgrInstance() (*UserMgr) {
	userMgrOnce.Do(func() {
		userMgrInstance = new(UserMgr)
	})

	return userMgrInstance
}

// 从数据库中加载数据
func (um *UserMgr) LoadData() ([]model.User, error) {
	sql := `SELECT uid, username FROM bs_user WHERE active`
	pg := postgres.GetInstance()
	dataSet, err := pg.Query(sql)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for _, rec := range dataSet {
		var user model.User
		uuid := rec["uid"].([]uint8)
		user.ID = string(uuid[:])
		user.Username = rec["username"].(string)
		users = append(users, user)
	}

	return users, nil
}

func (um *UserMgr) GetUser(username string) (model.User, error) {
	var user model.User

	sql := `SELECT uid, username FROM bs_user WHERE active AND username = $1`
	pg := postgres.GetInstance()
	dataSet, err := pg.Query(sql, username)
	if err != nil {
		return user, err
	}
	if len(dataSet) != 1 {
		err = errors.New("用户不存在、未激活或者表数据错误！")
		return user, err
	}

	rec := dataSet[0]
	uuid := rec["uid"].([]uint8)
	user.ID = string(uuid[:])
	user.Username = rec["username"].(string)

	return user, nil
}
