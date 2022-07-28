package entity

// UserDO mysql映射对象
type UserDO struct {
	Id          int    `db:"id"`
	Username    string `db:"username"`
	Nickname    string `db:"nickname"`
	Salt        string `db:"salt"`
	Password    string `db:"password"`
	ProfilePath string `db:"profile_path"`
	GmtCreate   string `db:"gmt_create"`
	GmtModified string `db:"gmt_modified"`
}
