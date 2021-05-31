package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mysql" // 匿名导入
)

var db *sql.DB

func init() {
	// 通过sql.Open来创建协程安全的sql.DB对象
	db, _ = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test1?charset=utf8")
	db.SetMaxOpenConns(1000) // 最大的连接个数，可自由设置，先设为1000
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

// DBConn: 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

// ParseRows: 将读到的sql语句信息转换为map类型的数组
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	// func (rs *Rows) Columns() ([]string, error)  Columns返回列名。如果Rows已经关闭会返回错误。
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
