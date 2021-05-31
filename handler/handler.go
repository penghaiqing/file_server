package handler

// 实现上传下载的接口

import (
	"encoding/json"
	dblayer "file_server/db"
	"file_server/meta"
	"file_server/util"
	"fmt"
	"io" // 输入输出
	"io/ioutil"
	"net/http" // 协议包引入
	"os"
	"strconv"
	"time"
)

// 上传文件的接口
/*
ResponseWriter接口被HTTP处理器用于构造HTTP回复。
type ResponseWriter interface {
    // Header返回一个Header类型值，该值会被WriteHeader方法发送。
    // 在调用WriteHeader或Write方法后再改变该对象是没有意义的。
    Header() Header
    // WriteHeader该方法发送HTTP回复的头域和状态码。
    // 如果没有被显式调用，第一次调用Write时会触发隐式调用WriteHeader(http.StatusOK)
    // WriterHeader的显式调用主要用于发送错误码。
    WriteHeader(int)
    // Write向连接中写入作为HTTP的一部分回复的数据。
    // 如果被调用时还未调用WriteHeader，本方法会先调用WriteHeader(http.StatusOK)
    // 如果Header中没有"Content-Type"键，
    // 本方法会使用包函数DetectContentType检查数据的前512字节，将返回值作为该键的值。
    Write([]byte) (int, error)
}
*/
func UploadHandler(w http.ResponseWriter, r *http.Request) { // w: 向用户返回数据的responseWriter对象 r：为接受用户请求的对象指针
	if r.Method == "GET" {
		// 返回上传的html页面
		data, err := ioutil.ReadFile("./static/view/index.html") // 通过读取文件来加载view中已经写好的html页面，使用相对路径，data为读到的文件内容
		if err != nil {                                          // 如果读取文件没有成功
			io.WriteString(w, "Internal Server Error")
			// func WriteString(w Writer, s string) (n int, err error)
			// WriteString函数将字符串s的内容写入w中。如果w已经实现了WriteString方法，函数会直接调用该方法。
			return
		}

		io.WriteString(w, string(data)) // 读取文档成功后直接将数据返回

	} else if r.Method == "POST" {
		//接受文件流并存储到本地的目录中
		file, head, err := r.FormFile("file") // 接受文件流,FormFile返回以key为键查询r.MultipartForm字段得到结果中的第一个文件和它的信息。
		if err != nil {
			fmt.Printf("Failed to get the data, err:%s\n", err.Error())
		}
		defer file.Close() // 关闭文件流

		// 此时需要要创建一个本地的文件来接收这个文件流
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			//Location: "D:\\Download_Chrome\\" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15;04;05"),
		}

		newFile, err := os.Create(fileMeta.Location) // Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）。
		//如果成功，返回的文件对象可用于I/O；对应的文件描述符具有O_RDWR模式。如果出错，错误底层类型是*PathError
		if err != nil {
			fmt.Printf("Failed to create file, err:%s", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file) // 通过copy函数，由newfile 接受file 中传过来的数据
		/*
			func Copy(dst Writer, src Reader) (written int64, err error)
			将src的数据拷贝到dst，直到在src上到达EOF或发生错误。返回拷贝的字节数和遇到的第一个错误。
			对成功的调用，返回值err为nil而非EOF，因为Copy定义为从src读取直到EOF，它不会将读取到EOF视为应报告的错误。
			如果src实现了WriterTo接口，本函数会调用src.WriteTo(dst)进行拷贝；
			否则如果dst实现了ReaderFrom接口，本函数会调用dst.ReadFrom(src)进行拷贝。
		*/
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s", err.Error())
			return
		}

		/*
			func (f *File) Seek(offset int64, whence int) (ret int64, err error)
			Seek设置下一次读/写的位置。
			offset为相对偏移量，而whence决定相对位置：0为相对文件开头，1为相对当前位置，2为相对文件结尾。
			它返回新的偏移量（相对开头）和可能的错误。
		*/
		newFile.Seek(0, 0) // 将newFile文件的seek位置移到0的位置

		// util 中的方法主要是通过哈希值来记录每个文件的信息，用于之后的查询
		fileMeta.FileSha1 = util.FileSha1(newFile) // 通过util的方法返回文件的 FileSha1，并赋值给fileMeta
		//meta.UpdateFileMeta(fileMeta) // 更新fileMeta的信息
		// 增加了mysql数据库持久化后，这里选择更新数据库
		meta.UpdateFileMetaDB(fileMeta)

		// TODO: 更新用户文件表()
		r.ParseForm() // 解析获得 username
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
			fileMeta.FileName, fileMeta.FileSize)
		if suc { // 更新成功
			http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
		} else {
			w.Write([]byte("Upload Failed."))
		}

		/*
			func Redirect(w ResponseWriter, r *Request, urlStr string, code int)
			Redirect回复请求一个重定向地址urlStr和状态码code。该重定向地址可以是相对于请求r的相对地址。
		*/

	}
}

// UploadSucHandler: 上传已完成 ，显示上传成功的信息
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler：获取文件的元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	// 增加mysql后的get方法
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if fMeta != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	} else {
		w.Write([]byte(`{"code":-1, "msg":"no such file"}`))
	}

}

// FileQuryHandler ：查询批量的文件元信息
func FileQuryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

// DownloadHandler: 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f) // 因为文件较小，所以选择使用ioutil中的方法直接读入到内存，大文件则需要设置流来读取
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 加http的响应头，让浏览器能够识别出下载
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileName+"\"")

	// 如果都成功了，直接将data信息返回到客户端即可(但是此时是 浏览器测试，则需要对下载进行页面的响应，需要上面几行代码进行处理)
	w.Write(data)

}

// FileMetaUpdateHandler: 更新元信息接口（重命名）
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	opType := r.Form.Get("op")
	filesha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	// opType 支持的操作类型
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden) // 不为0，则返回 403 的客户端错误
		return
	}
	// 如果方法不是 POST 方法，返回 405 客户端错误
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(filesha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 服务器错误 500
		return
	}
	w.WriteHeader(http.StatusOK) // statusOk = 200
	w.Write(data)
}

// FileDeleteHandler: 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filesha1 := r.Form.Get("filehash")

	// 此时需要对文件进行物理上的删除
	fMeta := meta.GetFileMeta(filesha1)
	os.Remove(fMeta.Location) // 此处可能会删除不成功报错，但是我们默认执行完后删除记录信息就表示删除了文件

	meta.RemoveFileMeta(filesha1)

	w.WriteHeader(http.StatusOK)
}
