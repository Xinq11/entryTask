package constant

type ErrCode int

const (
	ServerError ErrCode = 0 + iota
	InvalidSessionError
	UserExistedError
	UserNotExistError
	PasswordError
	InvalidParamsError
	Success
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
		return "UnKnow"
	}
}
