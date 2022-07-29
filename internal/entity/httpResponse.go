package entity

import (
	"EntryTask/constant"
	"fmt"
)

type HttpResponse struct {
	ErrCode constant.ErrCode `json:"errCode"`
	ErrMsg  string           `json:"errMsg"`
	Data    interface{}      `json:"data"`
}

func (hres HttpResponse) ToString() string {
	return fmt.Sprintf("errCode is %v, errMsg is %v, data is %v", hres.ErrCode, hres.ErrMsg, hres.Data)
}
