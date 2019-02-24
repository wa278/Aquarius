package lib

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var goDeployEnvironmentFlag string
//获取环境变量
func GetEnvFlag() string {
	if goDeployEnvironmentFlag == "" {
		return goDeployEnvironmentFlag
	}
	goDeployEnvironmentFlag = os.Getenv("env")
	if goDeployEnvironmentFlag == "" {
		goDeployEnvironmentFlag = "product"
		fmt.Printf("WARN", "Evironment Variable  not found, setting to `%s`.", goDeployEnvironmentFlag)
	}
	return goDeployEnvironmentFlag
}
//判断list是否包含str
func InListString(str string, list []string) bool {
	if len(str) > 0 {
		for _, itm := range list {
			if itm == str {
				return true
			}
		}
	}
	return false
}
//按分隔符将string 转化为list
func StringToList(str string,delimiter string) []string {
	var list = make([]string, 0)
	if str == "" {
		return list
	}
	for _, v := range strings.Split(str, delimiter) {
		v = strings.TrimSpace(v)
		if v != "" {
			list = append(list, v)
		}
	}
	return list
}

//调用系统方法获取goroutine id
func GetGoRoutineId() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	arr := strings.Fields(string(buf[:n]))
	goid := arr[1]
	return goid
}

// id 设计为系统时间+随机数
func GetUUID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano()) + fmt.Sprintf("%x", rand.Intn(10000))
}
//获取客户端Ip
func GetClientIP(r *http.Request) string {
	clientIP := r.RemoteAddr
	if strings.Index(clientIP, ":") > 0 {
		arr := strings.Split(clientIP, ":")
		clientIP = arr[0]
	}
	return clientIP
}
//将字符串转换成指针
func StringToPtr(str string) *string  {
	newStr := &str
	return newStr
}