package db

import (
	"fmt"
	mydb "file_server/db/mysql"
)
// UserSignup：通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`, `user_pwd`) values(?,?)")
		// insert ignore 会忽略数据库中已经存在的数据，如果数据库中没有数据，就直接插入新的数据；
		// 如果有数据的话就跳过这条数据，这样可以保留数据库中已经存在的数据，达到在间隙中插入数据的目的。

	if err != nil {
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}
	defer stmt.Close()

	// 执行 sql语句
	/*
		func (tx *Tx) Exec(query string, args ...interface{}) (Result, error)
		Exec执行命令，但不返回结果。例如执行insert和update。
		其中返回值 Result的定义：Result 是对已执行的 SQL 命令的总结
		type Result interface {
		// LastInsertId返回一个数据库生成的回应命令的整数。
		// 当插入新行时，一般来自一个"自增"列。
		// 不是所有的数据库都支持该功能，该状态的语法也各有不同。
		LastInsertId() (int64, error)

		// RowsAffected返回被update、insert或delete命令影响的行数。
		// 不是所有的数据库都支持该功能。
		RowsAffected() (int64, error)
	}
	*/
	ret, err := stmt.Exec(username, passwd)
	if err != nil{
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}

	// 因为sql中使用了 ignore 插入，需要判断下是否因有重复的key而没有进行数据的插入
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0{
		return true
	}
	return false
}
// UserSingin: 通过传入的用户名和加密密码 在数据库中进行比对判断是否合法
func UserSingin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	/*
	func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	Query执行一次查询，返回多行结果（即Rows），一般用于执行select命令。
	参数args表示query中的占位参数。
	Example:
	age := 27
	rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 返回值 Rows：
	// type Rows struct {
    // 		内含隐藏或非导出字段
	// }
	//是查询的结果。它的游标指向结果集的第零行，使用Next方法来遍历各行结果：
	

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, age)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	*/ 
	rows, err := stmt.Query(username)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}else if rows == nil{
		fmt.Println("username not found: " + username)
		return false
	}
}