package service

import (
	"EntryTask/constant"
	"EntryTask/database"
	"EntryTask/internal/entity"
	"testing"
)

// 注册
func TestUserService_SignUp(t *testing.T) {
	userService := &UserService{}
	database.MysqlInit()
	database.RedisInit()
	param := entity.UserDTO{
		Username: "ll",
		Password: "1234567",
	}
	res := userService.SignUp(param)
	if res.ErrCode != constant.Success {
		t.Error("signUp " + res.ErrCode.GetErrMsgByCode())
	} else {
		t.Log("signUp " + res.ErrCode.GetErrMsgByCode())
	}
}

// 登录
func TestUserService_SignIn(t *testing.T) {
	userService := &UserService{}
	database.MysqlInit()
	database.RedisInit()
	param := entity.UserDTO{
		Username: "xq",
		Password: "123456",
	}
	res := userService.SignIn(param)
	if res.ErrCode != constant.Success {
		t.Error("signIn " + res.ErrCode.GetErrMsgByCode())
	} else {
		t.Log("signIn " + res.ErrCode.GetErrMsgByCode())
	}
}

// 登出
func TestUserService_SignOut(t *testing.T) {
	userService := &UserService{}
	database.MysqlInit()
	database.RedisInit()
	param := entity.UserDTO{
		SessionID: "1ab00e57-10ae-4742-a461-fad16b536515",
	}
	res := userService.SignOut(param)
	if res.ErrCode != constant.Success {
		t.Error("signOut " + res.ErrCode.GetErrMsgByCode())
	} else {
		t.Log("signOut " + res.ErrCode.GetErrMsgByCode())
	}
}

// 查看用户信息
func TestUserService_GetUserInfo(t *testing.T) {
	userService := &UserService{}
	database.MysqlInit()
	database.RedisInit()
	param := entity.UserDTO{
		Username: "xq",
	}
	res := userService.GetUserInfo(param)
	if res.ErrCode != constant.Success {
		t.Error("getUserInfo " + res.ErrCode.GetErrMsgByCode())
	} else {
		t.Log("getUserInfo " + res.ErrCode.GetErrMsgByCode())
	}
}

// 编辑头像
func TestUserService_UpdateProfilePic(t *testing.T) {
	userService := &UserService{}
	database.MysqlInit()
	database.RedisInit()
	param := entity.UserDTO{
		Username:    "xq",
		ProfilePath: "ysds.jpg",
	}
	res := userService.UpdateProfilePic(param)
	if res.ErrCode != constant.Success {
		t.Error("updateNickName " + res.ErrCode.GetErrMsgByCode())
	} else {
		t.Log("updateNickName " + res.ErrCode.GetErrMsgByCode())
	}
}

// 编辑昵称
func TestUserService_UpdateNickName(t *testing.T) {
	userService := &UserService{}
	database.MysqlInit()
	database.RedisInit()
	param := entity.UserDTO{
		Username: "xq",
		Nickname: "young",
	}
	res := userService.UpdateNickName(param)
	if res.ErrCode != constant.Success {
		t.Error("updateNickName " + res.ErrCode.GetErrMsgByCode())
	} else {
		t.Log("updateNickName " + res.ErrCode.GetErrMsgByCode())
	}
}
