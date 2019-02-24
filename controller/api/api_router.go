package api

import (
	"fmt"
	"sync"
)

var (
	apiRouteTableLock sync.RWMutex
	apiRouteTable     = make(map[string]interface{})
)

func getApiAction(action string) (interface{}, error) {
	apiRouteTableLock.RLock()
	defer apiRouteTableLock.RUnlock()
	if obj, ok := apiRouteTable[action]; ok {
		return obj, nil
	} else {
		return nil, fmt.Errorf("The Action(%s) provided is invalid!", action)
	}
}

func registerApiAction(actionKey string, actionStruct interface{}) {
	apiRouteTableLock.Lock()
	defer apiRouteTableLock.Unlock()
	if _, ok := apiRouteTable[actionKey]; !ok {
		apiRouteTable[actionKey] = actionStruct
	}
}
