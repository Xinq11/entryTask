package manager

import (
	"EntryTask/database"
	"EntryTask/internal/entity"
	"time"
)

// 读取缓存
func GetUserInfoFromRedis(username string) (map[string]string, error) {
	cmd := database.RedisDB.HGetAll(username)
	return cmd.Val(), cmd.Err()
}

// 缓存userInfo
func CacheUserInfo(user entity.UserDO) error {
	userMap := make(map[string]interface{})
	userMap["username"] = user.Username
	if user.Nickname != "" {
		userMap["nickname"] = user.Nickname
	}
	if user.ProfilePath != "" {
		userMap["profilePath"] = user.ProfilePath
	}
	cmd := database.RedisDB.HMSet(user.Username, userMap)
	database.RedisDB.Expire(user.Username, 30*time.Minute)
	return cmd.Err()
}

// 删除缓存
func DelUserInfo(username string) error {
	cmd := database.RedisDB.HDel(username, "username", "nickname", "profilePath")
	return cmd.Err()
}
