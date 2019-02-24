
package web

import (
	"Aquarius/lib"
	"fmt"
	"net/http"
	"strings"
)
//获取配置模块
func getModules() []string {
	modules,err := lib.GetConfigStringList("Modules", []string{})
	if err!=nil{
		fmt.Printf(err.Error())
	}
	return modules
}

//获取模块名
func GetModuleName(r *http.Request) string {
	parts := getURLParts(r)
	if len(parts) >0{
		moduleName := strings.ToLower(parts[0])
		fmt.Println("modulename is " +moduleName)
		if lib.InListString(moduleName, getModules()) {
			return moduleName
		}
	}
	return "invalid"
}
//url分块
func getURLParts(r *http.Request) []string {
	return lib.StringToList(r.URL.Path, "/")
}

func getURLPartsWithoutModule(r *http.Request) []string {
	parts := getURLParts(r)
	result := make([]string, 0)
	for i, v := range parts {
		if i != 0 {
			result = append(result, v)
		}
	}
	return result
}
//获取controller
func GetControllerName(r *http.Request) string {
	parts := getURLPartsWithoutModule(r)
	ctrl := "gate"
	if len(parts) >= 1 {
		ctrl = strings.ToLower(parts[0])
	}
	return ctrl
}
//获取method
func GetMethodName(r *http.Request) string {
	parts := getURLPartsWithoutModule(r)
	method := "gate"
	if len(parts) >= 2 {
		method = parts[1]
	}
	return method
}