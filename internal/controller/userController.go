package controller

import (
	"EntryTask/constant"
	"EntryTask/internal/entity"
	"EntryTask/rpc/client"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strings"
	"time"
)

// 注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	setHttpHeader(w)
	// 参数校验
	/*
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
	*/
	var req entity.HttpRequest
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &req)
	if err != nil || len(req.Username) < 4 || len(req.Username) > 13 || !VerifyPassword(req.Password) {
		res := entity.HttpResponse{
			ErrCode: constant.InvalidParamsError,
			ErrMsg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:    "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		Username: req.Username,
		Password: req.Password,
	}
	rpcResponse := client.Client.Call("UserService.SignUp", userDTO)
	// 处理结果
	res := entity.HttpResponse{
		ErrCode: rpcResponse.ErrCode,
		ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
		Data:    "",
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 登录
func SignInHandler(w http.ResponseWriter, r *http.Request) {

	// 参数校验
	//username := r.PostFormValue("username")
	//password := r.PostFormValue("password")
	setHttpHeader(w)
	var req entity.HttpRequest
	body, err := ioutil.ReadAll(r.Body)
	logrus.Infoln(body)
	err = json.Unmarshal(body, &req)
	logrus.Infoln("controller receive param: ", req.Username, req.Password)
	if err != nil || len(req.Username) < 4 || len(req.Username) > 13 || !VerifyPassword(req.Password) {
		res := entity.HttpResponse{
			ErrCode: constant.InvalidParamsError,
			ErrMsg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:    "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		Username: req.Username,
		Password: req.Password,
	}

	rpcResponse := client.Client.Call("UserService.SignIn", userDTO)
	logrus.Infoln("controller receive rpc res: ", rpcResponse)
	var sessionID string
	// 处理结果
	if rpcResponse.ErrCode == constant.Success {
		dto := rpcResponse.Data.(entity.UserDTO)
		cookie := http.Cookie{
			Name:     "sessionID",
			Value:    dto.SessionID,
			Secure:   false,
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)
		sessionID = dto.SessionID
	}
	res := entity.HttpResponse{
		ErrCode: rpcResponse.ErrCode,
		ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
		Data:    sessionID,
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 登出
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	setHttpHeader(w)
	// 参数校验
	c, err := r.Cookie("sessionID")
	if err != nil || c == nil {
		res := entity.HttpResponse{
			ErrCode: constant.InvalidSessionError,
			ErrMsg:  constant.InvalidSessionError.GetErrMsgByCode(),
			Data:    "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		SessionID: c.Value,
	}
	rpcResponse := client.Client.Call("UserService.SignOut", userDTO)
	// 处理结果
	res := entity.HttpResponse{
		ErrCode: rpcResponse.ErrCode,
		ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
		Data:    "",
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 获取用户信息
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {

	setHttpHeader(w)
	// 参数校验
	c, err := r.Cookie("sessionID")
	if err != nil || c == nil {
		logrus.Infoln(err.Error())
		res := entity.HttpResponse{
			ErrCode: constant.InvalidSessionError,
			ErrMsg:  constant.InvalidSessionError.GetErrMsgByCode(),
			Data:    "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	logrus.Infoln("controller receive params: ", c.Value)
	// RPC
	userDTO := entity.UserDTO{
		SessionID: c.Value,
	}
	rpcResponse := client.Client.Call("UserService.GetUserInfo", userDTO)
	logrus.Infoln("controller receive rpc res: ", rpcResponse)
	// 处理结果
	res := entity.HttpResponse{
		ErrCode: rpcResponse.ErrCode,
		ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
	}
	if rpcResponse.ErrCode == constant.Success {
		dto := rpcResponse.Data.(entity.UserDTO)
		res.Data = entity.UserVO{
			Username:    dto.Username,
			Nickname:    dto.Nickname,
			ProfilePath: dto.ProfilePath,
		}
	} else {
		res.Data = ""
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 更新头像
func UpdateProfilePicHandler(w http.ResponseWriter, r *http.Request) {
	setHttpHeader(w)
	// 参数校验
	c, err := r.Cookie("sessionID")
	username := r.PostFormValue("username")
	logrus.Infoln("controller receive params: ", c.Value, username)
	if err != nil || c == nil || len(username) < 4 || len(username) > 13 {
		res := entity.HttpResponse{
			ErrCode: constant.InvalidParamsError,
			ErrMsg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:    "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}

	// 处理图片
	filePath, err := saveProfilePic(r, username)
	if err != nil {
		logrus.Error("userController.SignUpHandler saveProfilePic error: ", err.Error())
		res := entity.HttpResponse{
			ErrCode: constant.ServerError,
			ErrMsg:  constant.ServerError.GetErrMsgByCode(),
			Data:    "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		SessionID:   c.Value,
		ProfilePath: filePath,
	}
	rpcResponse := client.Client.Call("UserService.UpdateProfilePic", userDTO)
	// 处理结果
	res := entity.HttpResponse{
		ErrCode: rpcResponse.ErrCode,
		ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
	}
	if rpcResponse.ErrCode == constant.Success {
		res.Data = entity.UserVO{
			ProfilePath: filePath,
		}
	} else {
		res.Data = ""
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 更新昵称
func UpdateNicknameHandler(w http.ResponseWriter, r *http.Request) {
	setHttpHeader(w)
	// 参数校验
	c, err := r.Cookie("sessionID")
	// nickname := r.PostFormValue("nickname")
	var req entity.HttpRequest
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &req)
	nickname := []rune(req.Nickname)
	if err != nil || c == nil || len(nickname) == 0 || len(nickname) > 8 {
		res := entity.HttpResponse{
			ErrCode: constant.InvalidParamsError,
			ErrMsg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:    "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		SessionID: c.Value,
		Nickname:  req.Nickname,
	}
	rpcResponse := client.Client.Call("UserService.UpdateNickName", userDTO)
	// 处理结果
	res := entity.HttpResponse{
		ErrCode: rpcResponse.ErrCode,
		ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
	}
	if rpcResponse.ErrCode == constant.Success {
		res.Data = entity.UserVO{
			Nickname: req.Nickname,
		}
	} else {
		res.Data = ""
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

func setHttpHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
}

// 图片处理
func saveProfilePic(r *http.Request, username string) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		return "", err
	}
	fileName := strings.Split(header.Filename, ".")
	// 拼接文件名
	filePath := username + "-" + time.Now().Format("2006-01-02") + "." + fileName[len(fileName)-1]
	logrus.Infoln("saveProfilePic : ", filePath)
	f, err := os.OpenFile("img/"+filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// 校验密码 不合规返回false
func VerifyPassword(password string) bool {
	// 正则
	myRegex := regexp.MustCompile("^([a-zA-Z]*[0-9]*){7,16}$")
	res := myRegex.MatchString(password)
	return res
}
