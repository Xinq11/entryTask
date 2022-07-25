package entity

import "EntryTask/constant"

type HttpResponse struct {
	ErrCode constant.ErrCode `json:"errCode"`
	ErrMsg  string           `json:"errMsg"`
	Data    interface{}      `json:"data"`
}
