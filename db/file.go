package db

import (
	"database/sql" // ql包提供了保证SQL或类SQL数据库的泛用接口。使用sql包时必须注入（至少）一个数据库驱动。
	mydb "file_server/db/mysql"
	"fmt"
)

// 此文件用来实现对数据库的操作，如增删改查

// OnFileUploadFinished: 文件上传完成，保存meta
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values(?,?,?,?,1)")
	/*
		func (db *DB) Prepare(query string) (*Stmt, error)
		Prepare创建一个准备好的状态用于之后的查询和命令。返回值可以同时执行多个查询和命令。
	*/

	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	// 接下来是真正的调用Exec() 方法来执行语句
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 { // 这种情况是，filehash已经被存入数据库了，此时为重复存入，返回warning信息
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

// 接下来通过函数获取数据库中的文件元信息
// 先定义一个 存放文件信息的结构体
type TableFile struct {
	FileHash     string
	FileName     sql.NullString
	FileSize     sql.NullInt64
	FileAddr     sql.NullString
	FileUploadAt sql.NullString
}

// GetFileMeta: 从mysql中获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, file_addr,create_at from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr, &tfile.FileUploadAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}

// GetFileMetaList : 从mysql批量获取文件元信息
func GetFileMetaList(limit int) ([]TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_addr, file_name, file_size from tbl_file where status=1 limit ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	cloumns, _ := rows.Columns() // Columns返回列名
	values := make([]sql.RawBytes, len(cloumns))
	var tfiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
		//  func (rs *Rows) Scan(dest ...interface{}) error   Scan将当前行各列结果填充进dest指定的各个值中。
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	fmt.Println("len of tfiles : ",len(tfiles))
	return tfiles, nil

}
