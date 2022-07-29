package entity

import "fmt"

type HttpRequest struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (hreq HttpRequest) ToString() string {
	return fmt.Sprintf("username is %v, nickname is %v, password is %v", hreq.Username, hreq.Nickname, hreq.Password)
}
