package api

import (
	"fmt"
)

type Test struct {
}

func init() {
	registerApiAction("Test", &Test{})
}
//空载测试
func (this *Test) BenchMark(rawData []byte) (interface{}, int64, error) {
	fmt.Println("in benchmark")
	data := "Hello,world!"
	return data, 0, nil
}


