package meta

import (
	mydb "file_server/db"
)

// FileMeta：文件元信息结构
type FileMeta struct{
	FileSha1 string // 唯一标识一个文件
	FileName string
	FileSize int64
	Location string
	UploadAt string // 时间戳
}

var fileMetas map[string]FileMeta  // 定义一个map，key值为FileSha1，value 为 FileMeta结构体

// 初始化，用make为fileMete分配空间
func init()  {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta：新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB: 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool{
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}


// GetFile：根据filessha1返回文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}
// GetFileMetaDB: 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta,error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil{
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
		UploadAt: tfile.FileUploadAt.String,
	}
	return fmeta, nil
}


// RemoveFileMeta: 删除元信息，如果是真实场景需要加锁判断，保证多线程情况下的安全性
func RemoveFileMeta(filesha1 string)  {
	delete(fileMetas, filesha1)
}