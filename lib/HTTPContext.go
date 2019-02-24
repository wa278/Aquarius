package lib

import (
	"net/http"
	"sync"
)
//http上下文
type HTTPContext struct {
	Id             string
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	ClientIP       string
}

var httpContextMap map[string]*HTTPContext
var ContextRWLock sync.RWMutex

func init()  {
	httpContextMap = make(map[string]*HTTPContext)
}
//获取http上下文
func GetHTTPContext() *HTTPContext  {
	goid := GetGoRoutineId()
	ContextRWLock.RLock()
	defer ContextRWLock.RUnlock()
	if obj,ok := httpContextMap[goid];ok{
		return obj
	} else {
		return nil
	}
}

//创建http上下文
func BuildHttpContext(w http.ResponseWriter, r *http.Request) {
	goid := GetGoRoutineId()
	uuid := GetUUID()
	clientIp := GetClientIP(r)
	ContextRWLock.Lock()
	defer ContextRWLock.Unlock()
	httpContextMap[goid] = &HTTPContext{
		Id:             uuid,
		Request:        r,
		ResponseWriter: w,
		ClientIP:       clientIp,
	}
}

//删除http上下文
func RemoveHttpContext() {
	goid := GetGoRoutineId()
	ContextRWLock.Lock()
	defer ContextRWLock.Unlock()
	delete(httpContextMap, goid)
}