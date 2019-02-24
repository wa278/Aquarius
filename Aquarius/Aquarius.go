package main

import (
	"Aquarius/Services"
	"flag"
	"fmt"
	"reflect"
)

var server = flag.String("server", "HTTPServer", "HTTPServer or TCPServer")
var port = flag.String("port", "8765", "port")

func main() {
	flag.Parse();
	service := &Services.GoServices{}
	in := make([]reflect.Value, 0)
	in = append(in, reflect.ValueOf(*port))
	//查找server
	f := reflect.ValueOf(service).MethodByName(*server)
	if !f.IsValid() {
		fmt.Printf("server: %s is invalid.\n", *server)
		service.GetServices()
	}
	f.Call(in)
}
