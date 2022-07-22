package mapper

import (
	"EntryTask/database"
	"EntryTask/internal/entity"
	"fmt"
	"testing"
)

func TestQueryUserInfoByUsername(t *testing.T) {
	database.MysqlInit()
	userDTO, err := QueryUserInfoByUsername("xq")
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(userDTO)
	}
}

func TestInsertUserInfo(t *testing.T) {
	database.MysqlInit()
	userDO := entity.UserDO{
		Username:    "llll",
		Nickname:    "curry",
		Password:    "1234567",
		ProfilePath: "",
	}
	res, err := InsertUserInfo(userDO)
	if res == 0 || err != nil {
		if err != nil {
			t.Error(err.Error())
		}
	} else {
		t.Log(res)
	}
}

func TestUpdateNickName(t *testing.T) {
	database.MysqlInit()
	res, err := UpdateNickName("nickyoung", "x")
	fmt.Println(res, err)
	if res == 0 || err != nil {
		if err != nil {
			t.Error(err.Error())
		}
	} else {
		t.Log(res)
	}
}

func TestUpdateProfliePath(t *testing.T) {
	database.MysqlInit()
	res, err := UpdateProfilePath("xxxxxxxxxxxxxxxxx.jpg", "xq")
	if res == 0 || err != nil {
		if err != nil {
			t.Error(err.Error())
		}
	} else {
		t.Log(res)
	}
}
