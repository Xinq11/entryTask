package mapper

import (
	"EntryTask/database"
	"EntryTask/internal/entity"
	"time"
)

func QueryUserInfoByUsername(username string) (entity.UserDO, error) {
	db := database.MySqlDB
	var userDO entity.UserDO
	err := db.QueryRow("select id,username,nickname,salt,password,profile_path from user where username = ?", username).Scan(&userDO.Id, &userDO.Username, &userDO.Nickname, &userDO.Salt, &userDO.Password, &userDO.ProfilePath)
	return userDO, err
}

func InsertUserInfo(user entity.UserDO) (int64, error) {
	db := database.MySqlDB
	res, err := db.Exec("insert into user (username,salt,password,nickname,profile_path,gmt_create,gmt_modified) values (?,?,?,?,?,?,?)", user.Username, user.Salt, user.Password, user.Nickname, user.ProfilePath, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func UpdateNickName(nickname string, username string) (int64, error) {
	db := database.MySqlDB
	res, err := db.Exec("update user set nickname = ? , gmt_modified = ? where username = ?", nickname, time.Now(), username)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func UpdateProfilePath(profilePath string, username string) (int64, error) {
	db := database.MySqlDB
	res, err := db.Exec("update user set profile_path = ? , gmt_modified = ? where username = ?", profilePath, time.Now(), username)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
