package handler

// 初始化分块上传
func InitialMultipartUploadHandle(w http.ResponsWriter, r *http.Request) {
	// 1、解析用户请求参数
	// 2、获得redis的一个连接
	// 3、生成分块上传的初始化信息
	// 4、将初始化信息写入到redis缓存
	// 5、将响应初始化数据返回到客户端
}