package lib

import (
	"Aquarius/Log"
	"fmt"
)
//获取panic并记录
func CatchPanic() {
	if r := recover(); r != nil {
		PanicLog(r)
	}
}

func PanicLog(r interface{}) string {
	if r == nil {
		return ""
	}
	errmsg := fmt.Sprintf("panic: %v\n", r)
	log.Logx.LogError(errmsg)
	return errmsg
}
