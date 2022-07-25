package service

import (
	"EntryTask/constant"
	"EntryTask/internal/entity"
	"EntryTask/internal/manager"
	"EntryTask/internal/mapper"
	"EntryTask/rpc/rpcEntity"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"math/rand"
	_ "net/http/pprof"
	"time"
)

type UserService struct{}

// 注册
func (u *UserService) SignUp(user entity.UserDTO) rpcEntity.RpcResponse {
	logrus.Infoln(user.Username)
	// 校验用户是否存在
	userDO, err := mapper.QueryUserInfoByUsername(user.Username)
	if (err != nil && err != sql.ErrNoRows) || userDO.Username != "" {
		if userDO.Username != "" {
			return rpcEntity.RpcResponse{
				ErrCode: constant.UserExistedError,
			}
		} else {
			logrus.Error("userService.SignUp queryUserInfoByUsername error: ", err.Error())
			return rpcEntity.RpcResponse{
				ErrCode: constant.ServerError,
			}
		}
	}
	// 密码加密
	salt := generateSalt()
	password := encryptPassword(salt, user.Password)
	// 封装DO对象
	userDO = entity.UserDO{
		Salt:        salt,
		Password:    password,
		Username:    user.Username,
		Nickname:    user.Nickname,
		ProfilePath: user.ProfilePath,
	}
	// 存入mysql
	res, err := mapper.InsertUserInfo(userDO)
	if res == 0 || err != nil {
		logrus.Error("userService.SignUp insertUserInfo error: ", err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
	}
}

// 登录
func (u *UserService) SignIn(user entity.UserDTO) rpcEntity.RpcResponse {
	// 校验查看用户是否存在
	userDO, err := mapper.QueryUserInfoByUsername(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return rpcEntity.RpcResponse{
				ErrCode: constant.UserNotExistError,
			}
		} else {
			logrus.Error("userService.SignIn queryUserInfoByUsername error: ", err.Error())
			return rpcEntity.RpcResponse{
				ErrCode: constant.ServerError,
			}
		}
	}
	// 密码加密 核对密码
	password := encryptPassword(userDO.Salt, user.Password)
	if password != userDO.Password {
		return rpcEntity.RpcResponse{
			ErrCode: constant.PasswordError,
		}
	}
	// 生成session
	sessionID, err := generateSession()
	if err != nil {
		logrus.Error("userService.SignIn generateSession error: ", err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
	}
	// 写入redis
	err = manager.CacheSession(sessionID, user.Username)
	if err != nil {
		logrus.Error("userService.SignIn cacheSession error: ", err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
	}
	err = manager.CacheUserInfo(userDO)
	if err != nil {
		logrus.Error("userService.SignIn cacheUserInfo error: ", err.Error())
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
		Data: entity.UserDTO{
			SessionID: sessionID,
		},
	}
}

// 登出
func (u *UserService) SignOut(user entity.UserDTO) rpcEntity.RpcResponse {
	// 验证session
	_, err := manager.GetSession(user.SessionID)
	// session不存在 用户未登录或登录过期
	if err != nil {
		return rpcEntity.RpcResponse{
			ErrCode: constant.InvalidSessionError,
		}
	}
	// 删除sessionID
	err = manager.DelSession(user.SessionID)
	if err != nil {
		logrus.Error("userService.SignOut delSession error: ", err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
	}
}

// 查看用户信息
func (u *UserService) GetUserInfo(user entity.UserDTO) rpcEntity.RpcResponse {
	startTime := time.Now()
	defer logrus.Infoln(time.Now().Sub(startTime))
	// 验证session
	username, err := manager.GetSession(user.SessionID)
	// session不存在 用户未登录或登录过期
	if err != nil {
		return rpcEntity.RpcResponse{
			ErrCode: constant.InvalidSessionError,
		}
	}
	// 查redis
	userMap, err := manager.GetUserInfoFromRedis(username)
	// hgetall 一个不存在key，会返回空的map{}，不会返回error
	if err == nil && len(userMap) != 0 {
		return rpcEntity.RpcResponse{
			ErrCode: constant.Success,
			Data: entity.UserDTO{
				Username:    userMap["username"],
				Nickname:    userMap["nickname"],
				ProfilePath: userMap["profilePath"],
			},
		}
	}
	// 查mysql
	userDTO, err := mapper.QueryUserInfoByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return rpcEntity.RpcResponse{
				ErrCode: constant.UserNotExistError,
			}
		} else {
			logrus.Error("userService.GetUserInfo queryUserInfoByUsername error: ", err.Error())
			return rpcEntity.RpcResponse{
				ErrCode: constant.ServerError,
			}
		}
	}
	// 写入redis
	err = manager.CacheUserInfo(userDTO)
	if err != nil {
		logrus.Error("userService.GetUserInfo cacheUserInfo error: ", err.Error())
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
		Data: entity.UserDTO{
			Username:    userDTO.Username,
			Nickname:    userDTO.Nickname,
			ProfilePath: userDTO.ProfilePath,
		},
	}
}

// 编辑头像
func (u *UserService) UpdateProfilePic(user entity.UserDTO) rpcEntity.RpcResponse {
	// 验证session
	username, err := manager.GetSession(user.SessionID)
	// session不存在 用户未登录或登录过期
	if err != nil {
		return rpcEntity.RpcResponse{
			ErrCode: constant.InvalidSessionError,
		}
	}
	// 更新mysql
	res, err := mapper.UpdateProfilePath(user.ProfilePath, username)
	if res == 0 || err != nil {
		logrus.Error("userService.UpdateProfilePic updateProfliePath error: ", err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
	}
	// 删除缓存
	retryDelUserInfo(username)
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
		Data: entity.UserDTO{
			ProfilePath: "sss",
		},
	}
}

// 编辑昵称
func (u *UserService) UpdateNickName(user entity.UserDTO) rpcEntity.RpcResponse {
	// 验证session
	username, err := manager.GetSession(user.SessionID)
	// session不存在 用户未登录或登录过期
	if err != nil {
		return rpcEntity.RpcResponse{
			ErrCode: constant.InvalidSessionError,
		}
	}
	// 更新mysql
	res, err := mapper.UpdateNickName(user.Nickname, username)
	if res == 0 || err != nil {
		logrus.Error("userService.UpdateNickName updateNickName error: ", err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
	}
	// 删除缓存
	retryDelUserInfo(username)
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
		Data:    user.Nickname,
	}
}

// 生成session
func generateSession() (string, error) {
	cmd, err := uuid.NewV4()
	return cmd.String(), err
}

// 生成随机字符串
func generateSalt() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := make([]rune, 4)
	for i := range salt {
		salt[i] = letters[rand.Intn(len(letters))]
	}
	return string(salt)
}

// 密码加密 salt方式
func encryptPassword(salt string, password string) string {
	m := md5.New()
	m.Write([]byte(password))
	password = hex.EncodeToString(m.Sum(nil))
	m.Write([]byte(password + salt))
	password = hex.EncodeToString(m.Sum(nil))
	return password
}

// 重试
func retryDelUserInfo(username string) {
	err := errors.New("newError")
	retryTimes := 0
	for err != nil && retryTimes < 5 {
		err = manager.DelUserInfo(username)
		retryTimes++
	}
	if err != nil {
		logrus.Warn("retryDelUserInfo error: ", err.Error())
	}
}
