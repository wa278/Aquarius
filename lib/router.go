
package lib

import (
	"fmt"
	"sync"
)

var (

	routeTablesLock sync.RWMutex

	routeTables = make(map[string]map[string]interface{}) //路由表
)
//注册路由表
func RegisterObject(classType string, className string, object interface{}) {
	routeTablesLock.Lock()
	defer routeTablesLock.Unlock()
	if _, ok := routeTables[classType]; !ok {
		routeTables[classType] = make(map[string]interface{})
	}

	if _, ok := routeTables[classType][className]; !ok {
		routeTables[classType][className] = object
	}

}
//获取路由表中的obj
func GetObject(classType string, className string) (interface{}, error) {
	routeTablesLock.RLock()
	defer routeTablesLock.RUnlock()

	if objs, ok := routeTables[classType]; !ok {
		return nil, fmt.Errorf("router Table Type (%s) is undefined.", classType)
	} else {
		if obj, ok := objs[className]; ok {
			return obj, nil
		} else {
			return nil, fmt.Errorf(" Router Table Type (%s) Class Name (%s) is undefined.", classType, className)
		}
	}
}

