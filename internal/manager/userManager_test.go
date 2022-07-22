package manager

import (
	"EntryTask/database"
	"EntryTask/internal/entity"
	"testing"
)

func TestCacheUserInfo(t *testing.T) {
	database.RedisInit()
	userDO := entity.UserDO{
		Username: "xq",
		Nickname: "y",
	}
	err := CacheUserInfo(userDO)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log("success")
	}
}
