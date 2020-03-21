package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//SignupHandler:处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var (
		data        []byte
		err         error
		username    string
		password    string
		encPassword string
		succ        bool
	)
	if r.Method == "GET" {
		if data, err = ioutil.ReadFile("./static/view/signup.html"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username = r.Form.Get("username")
	password = r.Form.Get("password")
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}
	encPassword = util.Sha1([]byte(password + pwd_salt))
	succ = db.UserSignup(username, encPassword)
	if succ {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

//SignInHandler: 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var (
		pwdCheck    bool
		username    string
		password    string
		encPassword string
		token       string
		res         bool
		req         util.RespMsg
	)
	r.ParseForm()
	username = r.Form.Get("username")
	password = r.Form.Get("password")
	encPassword = util.Sha1([]byte(password + pwd_salt))
	// 1 校验用户名密码
	if pwdCheck = db.UserSignin(username, encPassword); !pwdCheck {
		w.Write([]byte("Failed"))
		return
	}
	// 2 生成访问凭证(tonken)
	token = GenToken(username)
	if res = db.UpdateToken(username, token); !res {
		w.Write([]byte("Failed"))
		return
	}
	// 3 登录成功 跳转到home.html
	//w.Write([]byte("http//" + r.Host + "/static/view/home.html"))
	req = util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			UserName string
			Token    string
		}{
			Location: "http//" + r.Host + "/static/view/home.html",
			UserName: username,
			Token:    token,
		},
	}
	w.Write(req.JSONBytes())
}

// GenToken 获取一个Token
func GenToken(username string) string {
	var (
		ts          string
		tokenPrefix string
	)
	// 40为字符串 md5(username+timestamp+token_salt)+timestamp[:8]
	ts = fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix = util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

//UserInfoHandler 获取用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		username string
		//token    string
		//isToken  bool
		user db.UserInfo
		err  error
		req  util.RespMsg
	)
	// 1 解析请求参数
	r.ParseForm()
	username = r.Form.Get("username")
	//token = r.Form.Get("token")
	// 2 验证token是否有效
	//if isToken = IsTokenValid(token); !isToken {
	//    w.WriteHeader(http.StatusForbidden)
	//    return
	//}
	// 3 查询用户信息
	if user, err = db.GetUserInfo(username); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4 组装并且相应用户数据
	req = util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(req.JSONBytes())
}

//IsTokenValid: 验证token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
