package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
	"github.com/dragonsun7/VCoinDataGather/config"
)

type Postgres struct {
	db *sql.DB
}

type Record map[string]interface{}
type DataSet []Record

var (
	sharedInstance *Postgres
	once           sync.Once
)

/* 单例 */
func GetInstance() (*Postgres) {
	once.Do(func() {
		sharedInstance = new(Postgres)

		cfg := config.GetInstance()
		driver := cfg.DB.Driver
		host := cfg.DB.Host
		port := cfg.DB.Port
		username := cfg.DB.Username
		password := cfg.DB.Password
		database := cfg.DB.Database

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, database)
		db, err := sql.Open(driver, dsn)
		if err != nil {
			panic(err)
		}
		sharedInstance.db = db
	})

	return sharedInstance
}

/* 关闭数据库(需要在适当的时候手动调用) */
func (pg *Postgres) CloseDB() {
	pg.db.Close()
}

func (pg *Postgres) Execute(sql string, args ...interface{}) (int64, error) {
	result, err := pg.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (pg *Postgres) Query(sql string, args ...interface{}) (DataSet, error) {
	rows, err := pg.db.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colCount := len(colNames)
	container := make([]interface{}, colCount)
	vars := make([]interface{}, colCount)
	for index := range container {
		container[index] = &vars[index]
	}

	dataSet := DataSet{}
	for rows.Next() {
		err = rows.Scan(container...)
		if err != nil {
			return nil, err
		}

		var rec = Record{}
		for index, colName := range colNames {
			rec[colName] = vars[index]
		}

		dataSet = append(dataSet, rec)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return dataSet, nil
}

func (pg *Postgres) BatchExecute(sql string, values [][]interface{}) (error) {
	return nil
	//tx, err := pg.db.Begin()
	//if err != nil {
	//	return err
	//}
	//
	//stmt, err := tx.Prepare(sql)
	//if err != nil {
	//	return err
	//}
	//defer stmt.Close()



	//stmt.Exec(args...)

	//tx, err := pg.DB.Begin()
	//stmt, err := tx.Prepare(sql)
	//result, err := stmt.Exec(1, 2, 3)
	//err := tx.Commit()

	//err := tx.Rollback()
}

func (pg *Postgres) GetOne() {

}

func (pg *Postgres) GetAll() {

}
