package constant

type ErrCode int

const (
	Success ErrCode = 0 + iota
	InvalidParamsError
	PasswordError
	UserNotExistError
	UserExistedError
	InvalidSessionError
	ServerError
)

func (code ErrCode) GetErrMsgByCode() string {
	switch code {
	case Success:
		return "success"
	case InvalidParamsError:
		return "InvalidParamsError"
	case PasswordError:
		return "PasswordError"
	case UserNotExistError:
		return "UserNotExistError"
	case UserExistedError:
		return "UserExistedError"
	case InvalidSessionError:
		return "InvalidSessionError"
	case ServerError:
		return "ServerError"
	default:
		return "Unknow"
	}
}
