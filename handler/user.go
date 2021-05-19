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
		//http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
		
		//Redirect 回复请求一个重定向地址urlStr和状态码code。该重定向地址可以是相对于请求r的相对地址。
	}

	r.ParseForm()  // 提取数据
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + pwd_salt)) 
	
	// 1、校验用户名和密码
	// 需要去 user 表中查询用户名 和 密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	fmt.Printf("pwdChecked: %t\n", pwdChecked)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		// fmt.Println("比较用户名密码后失败")
		return
	}
	// 到这里表示用户名和密码匹配正确
	// 2、生成访问凭证
	token := GenToken(username)
	// 然后要将这个token写入到数据库中去
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}
	// 3、登录成功后重定向到首页
	// w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	// 因为登录成功后再访问 其他的接口还是需要token信息才能成功的调用，所以在之后的传输中
	// 都要带上 token信息。建议使用 jason 格式来存放这些返回信息
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: struct{
			Location string
			Username string
			Token string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token: token,
		},
	}
	w.Write(resp.JSONBytes())
}
// UserInfoHandler: 查询并返回用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request)  {
	// 1、解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	// token := r.Form.Get("token")
	// // 2、验证token 是否有效, 定义一个函数去验证
	// isValdiToken := IsTokenValid(token)
	// if !isValdiToken{
	// 	w.WriteHeader(http.StatusForbidden) // 返回403 错误码
	// }
	// 3、查询用户信息
	// 此时应该在db的user.go 中增加查询的方法，在下面进行调用
	user, err := dblayer.GetUserInfo(username)
	if err != nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4、组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// GenToken：生成token
func GenToken(username string) string {
	// 40位字符： md5(username + timestamp + token_salt) + timestamp[:8]  , 32位md5 + 8 位时间
	ts := fmt.Sprintf("%x", time.Now().Unix())  // 格式化当前的时间戳
	// func Sprintf(format string, a ...interface{}) string
	// Sprintf根据format参数生成格式化的字符串并返回该字符串。
	tokenPrefix := util.MD5([]byte(username + ts + "_token_salt"))
	return tokenPrefix + ts[:8]
}
// IsTokenValid: 判断 token 是否一致
func IsTokenValid(token string) bool {
	// todo: 判断token 的时效性，是否过期
	// todo: 从数据库表tbl_user_token查询username对应的token信息
	// todo：对比两个token是否一致
	return true
}