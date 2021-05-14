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