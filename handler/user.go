package handler

import (
	"io/ioutil"
	"net/http"
	"file_server/util"
	dblayer "file_server/db"
)

const(
	pwd_salt = "**99"  // 用于Sha1() 去加密
)

// SignupHandler: 处理用户注册请求handler，是get方法就返回html文件，POST就对注册信息进行处理
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet{
		// 返回注册的html页面
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//io.WriteString(w, string(data)) // 读取文档成功后直接将数据返回
		w.Write(data)
		return
	}
	
	r.ParseForm() // 解析并返回 输入的用户信息
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	// 用户名及密码 合法性 判断
	if len(username)<3 || len(passwd)<5{
		w.Write([]byte("Invalid parameter"))
		return
	}
	// 对密码进行加密（使用 util 中的Sha1() 方法
	enc_passwd := util.Sha1([]byte(passwd + pwd_salt)) // 得到加密后的密码，存入数据库中才更安全
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc == true{
		w.Write([]byte("SUCCESS"))
		return
	}else{
		w.Write([]byte("Failed"))
	}
}
// SignInHandler: 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request)  {
	// 1、校验用户名和密码
	// 需要去 user 表中查询用户名 和 密码
	

	// 2、生成访问凭证


	// 3、登录成功后重定向到首页


}