package main

import (
	"Aquarius/lib"
	"Aquarius/protobuf"
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"time"
)

func main()  {
	client := http.Client{
		Timeout:3*time.Second,
	}
	req := &protobuf.ApiReq{
		Action:lib.StringToPtr("Test"),
		Method:lib.StringToPtr("BenchMark"),
		SystemId:lib.StringToPtr("123"),
	}
	data,_ := proto.Marshal(req)
	httpReq,_:= http.NewRequest(http.MethodPost,"http://127.0.0.1:8765/api/Gate",bytes.NewReader(data))
	resp,_:=client.Do(httpReq)
	fmt.Println(resp.StatusCode)
	responseBody,_ := ioutil.ReadAll(resp.Body)
	var response protobuf.ApiResp
	proto.Unmarshal(responseBody,&response)
	fmt.Println(response.GetCode())
	fmt.Println(response.GetErrmsg())
	fmt.Println(string(response.GetData()))

}


