package main

import (
	"file_server/handler"
	"fmt"
	"net/http"
)

func main() {
	// 静态资源目录映射
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	// 路由规则
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler)) // uplaod the files
	http.HandleFunc("/file/upload/suc", handler.HTTPInterceptor(handler.UploadSucHandler))

	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	/*
		func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
		HandleFunc注册一个处理器函数handler和对应的模式pattern（注册到DefaultServeMux）。
		ServeMux的文档解释了模式的匹配机制。
	*/

	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler))

	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	// 增加用户相关功能的路由规则
	http.HandleFunc("/", handler.SignInHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	// http.HandleFunc("/user/info", handler.UserInfoHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	/*
		func ListenAndServe
		func ListenAndServe(addr string, handler Handler) error
		ListenAndServe监听TCP地址addr，并且会使用handler参数调用Serve函数处理接收到的连接。
		handler参数一般会设为nil，此时会使用DefaultServeMux。
	*/
	err := http.ListenAndServe(":8080", nil) // 监听端口
	if err != nil {                          // 有错误返回错误信息
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}
