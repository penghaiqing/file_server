package meta

import (
	mydb "file_server/db"
	"sort"
)

// FileMeta：文件元信息结构
type FileMeta struct {
	FileSha1 string // 唯一标识一个文件
	FileName string
	FileSize int64
	Location string
	UploadAt string // 时间戳
}

var fileMetas map[string]FileMeta // 定义一个map，key值为FileSha1，value 为 FileMeta结构体

// 初始化，用make为fileMete分配空间
func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta：新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB: 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// GetFile：根据filessha1返回文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB: 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if tfile == nil || err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
		UploadAt: tfile.FileUploadAt.String,
	}
	return &fmeta, nil
}

// 获取批量的文件元信息列表
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

// 批量从mysql数据库中获取文件元信息
func GetLastFileMetasDB(limit int) ([]FileMeta, error) {
	tfiles, err := mydb.GetFileMetaList(limit)
	if err != nil {
		return make([]FileMeta, 0), err
	}
	tfilesm := make([]FileMeta, len(tfiles))
	for i := 0; i < len(tfilesm); i++ {
		tfilesm[i] = FileMeta{
			FileSha1: tfiles[i].FileHash,
			FileName: tfiles[i].FileName.String,
			FileSize: tfiles[i].FileSize.Int64,
			Location: tfiles[i].FileAddr.String,
			UploadAt: tfiles[i].FileUploadAt.String,
		}
	}
	return tfilesm, nil
}

// RemoveFileMeta: 删除元信息，如果是真实场景需要加锁判断，保证多线程情况下的安全性
func RemoveFileMeta(filesha1 string) {
	delete(fileMetas, filesha1)
}
