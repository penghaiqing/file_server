package main

import(
	"fmt"
	"net/http"
	"file_server/handler"
)

func main(){
	// 路由规则
	http.HandleFunc("/file/upload", handler.UploadHandler) // uplaod the files
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)

	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	/*
	func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
	HandleFunc注册一个处理器函数handler和对应的模式pattern（注册到DefaultServeMux）。
	ServeMux的文档解释了模式的匹配机制。
	*/

	/*
	func ListenAndServe
	func ListenAndServe(addr string, handler Handler) error
	ListenAndServe监听TCP地址addr，并且会使用handler参数调用Serve函数处理接收到的连接。
	handler参数一般会设为nil，此时会使用DefaultServeMux。
	*/
	
	http.HandleFunc("/file/download", handler.DownloadHandler)

	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	// 增加用户相关功能的路由规则
	//http.HandleFunc("/", handler.SignInHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)

	http.HandleFunc("/user/info", handler.UserInfoHandler)

	err := http.ListenAndServe(":8080",nil) // 监听端口
	if err != nil{ // 有错误返回错误信息
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}