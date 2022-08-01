package service

import (
	"EntryTask/config"
	"EntryTask/constant"
	"EntryTask/internal/entity"
	"EntryTask/internal/manager"
	"EntryTask/internal/mapper"
	"EntryTask/logger"
	"EntryTask/rpc/rpcEntity"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/gofrs/uuid"
	"math/rand"
	"time"
)

type UserService struct{}

// 注册
func (u *UserService) SignUp(user entity.UserDTO) (res rpcEntity.RpcResponse) {
	logger.Info("userService.SignUp receive userDTO is: " + user.ToString())
	defer func() {
		logger.Info("userService.SignUp response is: " + res.ToString())
	}()
	// 校验用户是否存在
	userDO, err := mapper.QueryUserInfoByUsername(user.Username)
	if (err != nil && err != sql.ErrNoRows) || userDO.Username != "" {
		logger.Error("userService.SignUp queryUserInfoByUsername error: " + err.Error())
		if userDO.Username != "" {
			return rpcEntity.RpcResponse{
				ErrCode: constant.UserExistedError,
			}
		} else {
			return rpcEntity.RpcResponse{
				ErrCode: constant.DataBaseError,
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
		Nickname:    config.DefaultNickname,
		ProfilePath: config.DefaultProfilePath,
	}
	// 存入mysql
	num, err := mapper.InsertUserInfo(userDO)
	if num == 0 || err != nil {
		logger.Error("userService.SignUp insertUserInfo error: " + err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.DataBaseError,
		}
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
	}
}

// 登录
func (u *UserService) SignIn(user entity.UserDTO) (res rpcEntity.RpcResponse) {
	logger.Info("userService.SignIn receive userDTO is: " + user.ToString())
	defer func() {
		logger.Info("userService.SignIn response is: " + res.ToString())
	}()
	// 校验查看用户是否存在
	userDO, err := mapper.QueryUserInfoByUsername(user.Username)
	if err != nil {
		logger.Error("userService.SignUp queryUserInfoByUsername error: " + err.Error())
		if err == sql.ErrNoRows {
			return rpcEntity.RpcResponse{
				ErrCode: constant.UserNotExistError,
			}
		} else {
			return rpcEntity.RpcResponse{
				ErrCode: constant.DataBaseError,
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
		logger.Error("userService.SignIn generateSession error: " + err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.ServerError,
		}
	}
	// 写入redis
	err = manager.CacheSession(sessionID, user.Username)
	if err != nil {
		logger.Error("userService.SignIn cacheSession error: " + err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.DataBaseError,
		}
	}
	if err = manager.CacheUserInfo(userDO); err != nil {
		logger.Error("userService.SignIn cacheUserInfo error: " + err.Error())
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
		Data: entity.UserDTO{
			SessionID: sessionID,
		},
	}
}

// 登出
func (u *UserService) SignOut(user entity.UserDTO) (res rpcEntity.RpcResponse) {
	logger.Info("userService.SignOut receive userDTO is: " + user.ToString())
	defer func() {
		logger.Info("userService.SignOut response is: " + res.ToString())
	}()
	// 删除sessionID
	if err := manager.DelSession(user.SessionID); err != nil {
		logger.Error("userService.SignOut delSession error: " + err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.DataBaseError,
		}
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
	}
}

// 查看用户信息
func (u *UserService) GetUserInfo(user entity.UserDTO) (res rpcEntity.RpcResponse) {
	logger.Info("userService.GetUserInfo receive userDTO is: " + user.ToString())
	defer func() {
		logger.Info("userService.GetUserInfo response is: " + res.ToString())
	}()
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
	userDO, err := mapper.QueryUserInfoByUsername(username)
	if err != nil {
		logger.Error("userService.GetUserInfo queryUserInfoByUsername error: " + err.Error())
		if err == sql.ErrNoRows {
			return rpcEntity.RpcResponse{
				ErrCode: constant.UserNotExistError,
			}
		} else {
			return rpcEntity.RpcResponse{
				ErrCode: constant.DataBaseError,
			}
		}
	}
	// 写入redis
	if err = manager.CacheUserInfo(userDO); err != nil {
		logger.Error("userService.GetUserInfo cacheUserInfo error: " + err.Error())
	}
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
		Data: entity.UserDTO{
			Username:    userDO.Username,
			Nickname:    userDO.Nickname,
			ProfilePath: userDO.ProfilePath,
		},
	}
}

// 编辑头像
func (u *UserService) UpdateProfilePic(user entity.UserDTO) (res rpcEntity.RpcResponse) {
	logger.Info("userService.UpdateProfilePic receive userDTO is: " + user.ToString())
	defer func() {
		logger.Info("userService.UpdateProfilePic response is: " + res.ToString())
	}()
	// 验证session
	username, err := manager.GetSession(user.SessionID)
	// session不存在 用户未登录或登录过期
	if err != nil {
		return rpcEntity.RpcResponse{
			ErrCode: constant.InvalidSessionError,
		}
	}
	// 更新mysql
	_, err = mapper.UpdateProfilePath(user.ProfilePath, username)
	if err != nil {
		logger.Error("userService.UpdateProfilePic updateProfilePath error: " + err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.DataBaseError,
		}
	}
	// 删除缓存
	go retryDelUserInfo(username)
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
	}
}

// 编辑昵称
func (u *UserService) UpdateNickName(user entity.UserDTO) (res rpcEntity.RpcResponse) {
	logger.Info("userService.UpdateNickName receive userDTO is: " + user.ToString())
	defer func() {
		logger.Info("userService.UpdateNickName response is: " + res.ToString())
	}()
	// 验证session
	username, err := manager.GetSession(user.SessionID)
	// session不存在 用户未登录或登录过期
	if err != nil {
		return rpcEntity.RpcResponse{
			ErrCode: constant.InvalidSessionError,
		}
	}
	// 更新mysql
	_, err = mapper.UpdateNickName(user.Nickname, username)
	if err != nil {
		logger.Error("userService.UpdateNickName updateNickName error: " + err.Error())
		return rpcEntity.RpcResponse{
			ErrCode: constant.DataBaseError,
		}
	}
	// 删除缓存
	go retryDelUserInfo(username)
	return rpcEntity.RpcResponse{
		ErrCode: constant.Success,
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
	var letters = []rune(config.Letters)
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
		logger.Warn("userService.retryDelUserInfo DelUserInfo error: " + err.Error())
	}
}
