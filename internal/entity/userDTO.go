package entity

// UserDTO RPC传输对象
type UserDTO struct {
	Username    string
	Nickname    string
	ProfilePath string
	Password    string
	SessionID   string
}
