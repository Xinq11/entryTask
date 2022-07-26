package entity

type UserVO struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	ProfilePath string `json:"profilePath"`
}
