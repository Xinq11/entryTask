package entity

import "fmt"

// UserDTO RPC传输对象
type UserDTO struct {
	Username    string
	Nickname    string
	ProfilePath string
	Password    string
	SessionID   string
}

func (userDTO UserDTO) ToString() string {
	return fmt.Sprintf("Username is %v, Nickname is %v, ProfilePath is %v, Password is %v, SessionID is %v",
		userDTO.Username, userDTO.Nickname, userDTO.ProfilePath, userDTO.Password, userDTO.SessionID)
}
