package api

import (
	"Aquarius/lib"
)

var controlName = "api"

func registerControllerName(ctrlName string, obj interface{}) {
	lib.RegisterObject("controller_"+controlName, "Gate", &Gate{})
}

func GetControllerInstance(ctrlName string) (interface{}, error) {
	return lib.GetObject("controller_"+controlName, "Gate")
}
