package mysql

import(
	"fmt"
	"os"
	"database/sql"
	_ "github.com/mysql" // 匿名导入
)

var db *sql.DB

func init(){
	// 通过sql.Open来创建协程安全的sql.DB对象
	db, _ = sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/test1?charset=utf8")
	db.SetMaxOpenConns(1000) // 最大的连接个数，可自由设置，先设为1000
	err := db.Ping()
	if err != nil{
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}
// DBConn: 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}