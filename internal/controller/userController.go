package controller

import (
	"EntryTask/constant"
	"EntryTask/internal/entity"
	"EntryTask/rpc/client"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// 注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 参数校验
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) == 0 || len(username) > 8 || !VerifyPassword(password) {
		res := entity.HttpResponse{
			Err_code: constant.InvalidParamsError,
			Err_msg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:     "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		Username: username,
		Password: password,
	}
	rpcResponse := client.Client.Call("UserService.SignUp", userDTO)
	// 处理结果
	res := entity.HttpResponse{
		Err_code: rpcResponse.Err_code,
		Err_msg:  rpcResponse.Err_code.GetErrMsgByCode(),
		Data:     "",
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 登录
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 参数校验
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) == 0 || len(username) > 8 || !VerifyPassword(password) {
		res := entity.HttpResponse{
			Err_code: constant.InvalidParamsError,
			Err_msg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:     "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		Username: username,
		Password: password,
	}
	rpcResponse := client.Client.Call("UserService.SignIn", userDTO)
	// 处理结果
	if rpcResponse.Err_code == constant.Success {
		dto := rpcResponse.Data.(entity.UserDTO)
		cookie := http.Cookie{
			Name:  "sessionID",
			Value: dto.SessionID,
		}
		http.SetCookie(w, &cookie)
	}
	res := entity.HttpResponse{
		Err_code: rpcResponse.Err_code,
		Err_msg:  rpcResponse.Err_code.GetErrMsgByCode(),
		Data:     "",
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 登出
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 参数校验
	c, err := r.Cookie("sessionID")
	if err != nil || c == nil {
		res := entity.HttpResponse{
			Err_code: constant.InvalidSessionError,
			Err_msg:  constant.InvalidSessionError.GetErrMsgByCode(),
			Data:     "",
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
		Err_code: rpcResponse.Err_code,
		Err_msg:  rpcResponse.Err_code.GetErrMsgByCode(),
		Data:     "",
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// 获取用户信息
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 参数校验
	c, err := r.Cookie("sessionID")
	if err != nil || c == nil {
		res := entity.HttpResponse{
			Err_code: constant.InvalidSessionError,
			Err_msg:  constant.InvalidSessionError.GetErrMsgByCode(),
			Data:     "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		SessionID: c.Value,
	}
	rpcResponse := client.Client.Call("UserService.GetUserInfo", userDTO)
	// 处理结果
	res := entity.HttpResponse{
		Err_code: rpcResponse.Err_code,
		Err_msg:  rpcResponse.Err_code.GetErrMsgByCode(),
	}
	if rpcResponse.Err_code == constant.Success {
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
	w.Header().Set("Content-Type", "application/json")
	// 参数校验
	c, err := r.Cookie("sessionID")
	username := r.PostFormValue("username")
	if err != nil || c == nil || len(username) == 0 || len(username) > 8 {
		res := entity.HttpResponse{
			Err_code: constant.InvalidParamsError,
			Err_msg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:     "",
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
			Err_code: constant.ServerError,
			Err_msg:  constant.ServerError.GetErrMsgByCode(),
			Data:     "",
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
		Err_code: rpcResponse.Err_code,
		Err_msg:  rpcResponse.Err_code.GetErrMsgByCode(),
	}
	if rpcResponse.Err_code == constant.Success {
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
	w.Header().Set("Content-Type", "application/json")
	// 参数校验
	c, err := r.Cookie("sessionID")
	nickname := r.PostFormValue("nickname")
	if err != nil || c == nil || len(nickname) == 0 || len(nickname) > 16 {
		res := entity.HttpResponse{
			Err_code: constant.InvalidParamsError,
			Err_msg:  constant.InvalidParamsError.GetErrMsgByCode(),
			Data:     "",
		}
		js, _ := json.Marshal(res)
		w.Write(js)
		return
	}
	// RPC
	userDTO := entity.UserDTO{
		SessionID: c.Value,
		Nickname:  nickname,
	}
	rpcResponse := client.Client.Call("UserService.UpdateNickName", userDTO)
	// 处理结果
	res := entity.HttpResponse{
		Err_code: rpcResponse.Err_code,
		Err_msg:  rpcResponse.Err_code.GetErrMsgByCode(),
		Data:     "",
	}
	js, _ := json.Marshal(res)
	w.Write(js)
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
	filePath := "img/" + username + "-" + time.Now().Format("2006-01-02-15:04:05") + "." + fileName[len(fileName)-1]
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
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
	myRegex := regexp.MustCompile("^[a-zA-Z][0-9]{7,16}$")
	res := myRegex.MatchString(password)
	return res
}
