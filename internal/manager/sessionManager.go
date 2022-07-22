package manager

import (
	"EntryTask/database"
	"time"
)

// 获取session
func GetSession(sessionID string) (string, error) {
	cmd := database.RedisDB.Get(sessionID)
	return cmd.Val(), cmd.Err()
}

// 删除缓存
func DelSession(sessionID string) error {
	cmd := database.RedisDB.Del(sessionID)
	return cmd.Err()
}

// 写入缓存
func CacheSession(sessionID string, username string) error {
	cmd := database.RedisDB.Set(sessionID, username, 30*time.Minute)
	return cmd.Err()
}
