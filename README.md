# file_server
a file_server with golang for practicing 

完成文件的上传和保存功能，其中对文件元信息的保存还没完善，需要后续保存到数据库中，才能保证其安全性。
简单的思路：
main.go 中调用 handler中的 UploadHandler 函数，根据用户的请求选择 GET 还是 POST 方法对文件流进行处理；GET是将准备好的静态页面返回，直接通过IO内置函数 WriteString() 将读到的html信息传到ResponseWriter对象即可；POST则是接收文件流并将其存储到本地目录中，最简单的处理是可以直接os.Create() 一个新文件，再调用io.Copy() 方法将文件流存储到新文件，此处我们为文件信息定义 FileMeta 这样一个 文件元信息结构 ，对文件流的信息进行存储，尤其是通过util/util.go 中的方法计算出每个文件的哈希标识，为之后对文件的查询修改操作提供便利。为了交互友好性，增加了UploadSucHandler函数，文件上传完成后，由main函数调用以显示上传成功的信息。
