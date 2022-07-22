package entity

import "EntryTask/constant"

type HttpResponse struct {
	Err_code constant.ErrCode
	Err_msg  string
	Data     interface{}
}
