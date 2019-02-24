package api

import (
	"Aquarius/Log"
	"Aquarius/lib"
	"Aquarius/web"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Gate struct {
}

var ApiLog  = log.NewLog("../log","api",4,10,1024*1024*100,log.LOG_SHIFT_BY_DAY)
func init() {
	//注册controller
	registerControllerName("Gate", &Gate{})
}
func (this *Gate) Gate(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Nanosecond()
	fmt.Println("in gate")
	defer lib.CatchPanic()
	rawData, rawDataErr := ioutil.ReadAll(r.Body)
	util := NewApiUtil()
	if rawDataErr != nil {
		retval := util.GetResult(EC_DEFAULT_ERR, rawDataErr.Error(), nil)
		web.Echo(w, retval)
		return
	}
	defer r.Body.Close()
	clientIP := lib.GetClientIP(r)
	//服务的统一入口
	retval, _ := util.Run(rawData, clientIP)
	web.Echo(w, retval)
	ApiLog.LogInfo("ip : %s's request cost %d ns",clientIP,time.Now().Nanosecond()-startTime)

}
