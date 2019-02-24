package api

import (
	"Aquarius/lib"
	"Aquarius/protobuf"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
	"strings"
)

const (
	EC_DEFAULT_ERR int64 = 1
)

type ApiUtil struct {
}

func NewApiUtil() *ApiUtil {
	return &ApiUtil{}
}
//api主流程
func (this *ApiUtil) Run(rawData []byte, clientIP string) (ret []byte, isLog bool) {
	fmt.Println("in apirun")
	maxPacketSize, err := lib.GetConfigInt64("ApiMaxAllowedPacket", 2)
	if err != nil {

	}
	//包大小
	if len(rawData) > 1024*1024*int(maxPacketSize) {
		return this.GetResult(EC_DEFAULT_ERR, fmt.Sprintf("The requested packet cannot exceed %d MB", maxPacketSize), nil), false
	}
	//解析
	var req = &protobuf.ApiReq{}
	if err := proto.Unmarshal(rawData, req); err != nil {
		return this.GetResult(EC_DEFAULT_ERR, "The received pb data cannot be parsed: "+err.Error(), nil), false
	}

	systemKey := req.GetSystemId()
	//鉴权
	isPrivPass, isPrivPassErr := this.ClientIpAuth(systemKey, clientIP)
	if isPrivPassErr != nil {
		return this.GetResult(EC_DEFAULT_ERR, isPrivPassErr.Error(), nil), true
	}
	if !isPrivPass {
		return this.GetResult(EC_DEFAULT_ERR, fmt.Sprintf("ip : %s, systemId: %s no authority", clientIP, systemKey), nil), true
	}
	//获取处理器
	apiActionObj, apiActionObjErr := getApiAction(req.GetAction())
	if apiActionObjErr != nil {
		return this.GetResult(EC_DEFAULT_ERR, apiActionObjErr.Error(), nil), true
	}

	f := reflect.ValueOf(apiActionObj).MethodByName(strings.Title(req.GetMethod()))
	if !f.IsValid() {
		return this.GetResult(EC_DEFAULT_ERR, fmt.Sprintf("Action : %s Method : %s is invalid.", req.GetAction(), req.GetMethod()), nil), true
	}

	method, _ := f.Interface().(func([]byte) (interface{}, int64, error))
	result, code, err := method(rawData)
	if err != nil {
		return this.GetResult(code, err.Error(), result), true
	} else {
		return this.GetResult(0, "", result), true
	}
}
//鉴权
func (this *ApiUtil) ClientIpAuth(systemKey string, clientIP string) (bool, error) {
	//TODO
	return true, nil
}
//封装返回信息
func (this *ApiUtil) GetResult(code int64, message string, data interface{}) []byte {
	fmt.Println("start get result")
	newcode := int32(code)
	resp := &protobuf.ApiResp{
		Code:   &newcode,
		Errmsg: &(message),
		Data:   []byte(data.(string)),
	}
	result, err := proto.Marshal(resp)
	if err!=nil{
		fmt.Println("get result wrong")
	}
	return result
}
