package web

import (
	"net/http"
)

//输出信息给用户
func Echo(w http.ResponseWriter, data []byte) {
	if _, err := w.Write(data); err != nil {
	}
}


