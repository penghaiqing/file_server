package handler

import (
	"io/ioutil"
	"net/http"
	"file_server/util"
	dblayer "file_server/db"
	"fmt"
	"time"
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
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		// http.Redirect(w, r, "./static/view/signin.html", http.StatusFound)
		// return
		
		//Redirect 回复请求一个重定向地址urlStr和状态码code。该重定向地址可以是相对于请求r的相对地址。
	}

	r.ParseForm()  // 提取数据
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + pwd_salt)) 
	
	// 1、校验用户名和密码
	// 需要去 user 表中查询用户名 和 密码
	pwdChecked := dblayer.UserSingin(username, encPasswd)
	fmt.Printf("pwdChecked: %t\n", pwdChecked)
	if !pwdChecked {
		w.Write([]byte("FAILED2"))
		// fmt.Println("比较用户名密码后失败")
		return
	}
	// 到这里表示用户名和密码匹配正确
	// 2、生成访问凭证
	token := GenToken(username)
	// 然后要将这个token写入到数据库中去
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED1"))
		return
	}
	// 3、登录成功后重定向到首页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

}

func GenToken(username string) string {
	// 40位字符： md5(username + timestamp + token_salt) + timestamp[:8]  , 32位md5 + 8 位时间
	ts := fmt.Sprintf("%x", time.Now().Unix())  // 格式化当前的时间戳
	tokenPrefix := util.MD5([]byte(username + ts + "_token_salt"))
	return tokenPrefix + ts[:8]
}