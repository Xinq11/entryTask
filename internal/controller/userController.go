package controller

import (
	"EntryTask/config"
	"EntryTask/constant"
	"EntryTask/internal/entity"
	"EntryTask/logger"
	rpcClient "EntryTask/rpc/client"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

// 注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		setHttpHeader(w)
		// 参数校验
		var req entity.HttpRequest
		err := paramParse(r, &req)
		logger.Info("userController.SignUpHandler receive HttpRequest is: " + req.ToString())
		if err != nil || len(req.Username) < 4 || len(req.Username) > 13 || len(req.Password) < 4 || len(req.Password) > 13 {
			res := entity.HttpResponse{
				ErrCode: constant.InvalidParamsError,
				ErrMsg:  constant.InvalidParamsError.GetErrMsgByCode(),
				Data:    "",
			}
			// json.Marshal参数类型不为无效类型(Chan,Func)或无效值(math.Inf,math.NaN)则可忽略err
			js, _ := json.Marshal(res)
			w.Write(js)
			return
		}
		// RPC
		userDTO := entity.UserDTO{
			Username: req.Username,
			Password: req.Password,
		}
		rpcResponse := rpcClient.Client.Call("UserService.SignUp", userDTO)
		// 处理结果
		res := entity.HttpResponse{
			ErrCode: rpcResponse.ErrCode,
			ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
			Data:    "",
		}
		logger.Info("userController.SignUpHandler response is: " + res.ToString())
		js, _ := json.Marshal(res)
		w.Write(js)
	}
}

// 登录
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		setHttpHeader(w)
		// 参数校验
		var req entity.HttpRequest
		err := paramParse(r, &req)
		logger.Info("userController.SignInHandler receive HttpRequest is: " + req.ToString())
		if err != nil || len(req.Username) < 4 || len(req.Username) > 13 || len(req.Password) < 4 || len(req.Password) > 13 {
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
		rpcResponse := rpcClient.Client.Call("UserService.SignIn", userDTO)
		// 处理结果
		if rpcResponse.ErrCode == constant.Success {
			dto := rpcResponse.Data.(entity.UserDTO)
			cookie := http.Cookie{
				Name:   "sessionID",
				Value:  dto.SessionID,
				MaxAge: 1800,
			}
			http.SetCookie(w, &cookie)
		}
		res := entity.HttpResponse{
			ErrCode: rpcResponse.ErrCode,
			ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
			Data:    "",
		}
		logger.Info("userController.SignInHandler response is: " + res.ToString())
		js, _ := json.Marshal(res)
		w.Write(js)
	}
}

// 登出
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		setHttpHeader(w)
		// 参数校验
		sessionID, err := r.Cookie("sessionID")
		if err != nil {
			logger.Error("userController.SignOutHandler getSessionID error: " + err.Error())
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
			SessionID: sessionID.Value,
		}
		rpcResponse := rpcClient.Client.Call("UserService.SignOut", userDTO)
		// 处理结果
		res := entity.HttpResponse{
			ErrCode: rpcResponse.ErrCode,
			ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
			Data:    "",
		}
		logger.Info("userController.SignOutHandler response is: " + res.ToString())
		js, _ := json.Marshal(res)
		w.Write(js)
	}
}

// 获取用户信息
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		setHttpHeader(w)
		// 参数校验
		sessionID, err := r.Cookie("sessionID")
		if err != nil {
			logger.Error("userController.GetUserInfoHandler getSessionID error: " + err.Error())
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
			SessionID: sessionID.Value,
		}
		rpcResponse := rpcClient.Client.Call("UserService.GetUserInfo", userDTO)
		// 处理结果
		res := entity.HttpResponse{
			ErrCode: rpcResponse.ErrCode,
			ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
			Data:    "",
		}
		if rpcResponse.ErrCode == constant.Success {
			dto := rpcResponse.Data.(entity.UserDTO)
			res.Data = entity.UserVO{
				Username:    dto.Username,
				Nickname:    dto.Nickname,
				ProfilePath: dto.ProfilePath,
			}
		}
		logger.Info("userController.GetUserInfoHandler response is: " + res.ToString())
		js, _ := json.Marshal(res)
		w.Write(js)
	}
}

// 更新头像
func UpdateProfilePicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		setHttpHeader(w)
		// 参数校验
		sessionID, err := r.Cookie("sessionID")
		username := r.PostFormValue("username")
		logger.Info("userController.UpdateProfilePicHandler receive username is: " + username)
		if err != nil || len(username) < 4 || len(username) > 13 {
			res := entity.HttpResponse{
				ErrCode: constant.InvalidParamsError,
				ErrMsg:  constant.InvalidParamsError.GetErrMsgByCode(),
				Data:    "",
			}
			js, _ := json.Marshal(res)
			w.Write(js)
			return
		}
		// 保存图片
		filePath, err := saveProfilePic(r, username)
		if err != nil {
			logger.Error("userController.UpdateProfilePicHandler saveProfilePic error: " + err.Error())
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
			SessionID:   sessionID.Value,
			ProfilePath: filePath,
		}
		rpcResponse := rpcClient.Client.Call("UserService.UpdateProfilePic", userDTO)
		// 处理结果
		res := entity.HttpResponse{
			ErrCode: rpcResponse.ErrCode,
			ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
			Data:    "",
		}
		if rpcResponse.ErrCode == constant.Success {
			res.Data = entity.UserVO{
				ProfilePath: filePath,
			}
		} else {
			go delProfilePic(filePath)
		}
		logger.Info("userController.UpdateProfilePicHandler response is: " + res.ToString())
		js, _ := json.Marshal(res)
		w.Write(js)
	}
}

// 更新昵称
func UpdateNicknameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		setHttpHeader(w)
		// 参数校验
		var req entity.HttpRequest
		sessionID, cookieErr := r.Cookie("sessionID")
		paramErr := paramParse(r, &req)
		logger.Info("userController.UpdateNicknameHandler receive HttpRequest is: " + req.ToString())
		if cookieErr != nil || paramErr != nil || len([]rune(req.Nickname)) < 1 || len([]rune(req.Nickname)) > 8 {
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
			SessionID: sessionID.Value,
			Nickname:  req.Nickname,
		}
		rpcResponse := rpcClient.Client.Call("UserService.UpdateNickName", userDTO)
		// 处理结果
		res := entity.HttpResponse{
			ErrCode: rpcResponse.ErrCode,
			ErrMsg:  rpcResponse.ErrCode.GetErrMsgByCode(),
			Data:    "",
		}
		if rpcResponse.ErrCode == constant.Success {
			res.Data = entity.UserVO{
				Nickname: req.Nickname,
			}
		}
		logger.Info("userController.UpdateNicknameHandler response is: " + res.ToString())
		js, _ := json.Marshal(res)
		w.Write(js)
	}
}

// 设置响应头
func setHttpHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
}

// 解析参数
func paramParse(r *http.Request, req *entity.HttpRequest) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, req); err != nil {
		return err
	}
	return nil
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
	f, err := os.OpenFile(config.FilePath+filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(f, file); err != nil {
		return "", err
	}
	return filePath, nil
}

// 删除文件
func delProfilePic(filePath string) {
	if err := os.Remove(config.FilePath + filePath); err != nil {
		logger.Error("userController.delProfilePic error: " + err.Error())
	}
}
