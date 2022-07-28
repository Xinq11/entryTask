package entity

// UserVO 前端映射对象
type UserVO struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	ProfilePath string `json:"profilePath"`
}
