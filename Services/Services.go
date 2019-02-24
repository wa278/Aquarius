package Services

import (
	"Aquarius/Log"
	"Aquarius/controller/api"
	"Aquarius/lib"
	"Aquarius/web"
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

type GoServices struct{}

func (this *GoServices) GetServices() {
	v := reflect.ValueOf(this)
	num := v.NumMethod();
	for i := 0; i < num; i++ {
		fname := v.Type().Method(i).Name
		fmt.Println(fname)
	}
}

func (this *GoServices) HTTPServer(port string) {
	fmt.Println("in httpserver")
	stopChan := make(chan os.Signal)
	//监听系统signal
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	//设置系统路由
	mux := http.NewServeMux()
	//监控系统信息
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	//以根目录作为系统的统一入口
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer lib.CatchPanic()
		lib.BuildHttpContext(w, r)
		defer lib.RemoveHttpContext()
		ctrlName := web.GetControllerName(r)

		modName := web.GetModuleName(r)

		var ctrlObj interface{}
		var ctrlObjErr error
		switch modName {

		case "api":
			fmt.Println("in api")
			ctrlObj, ctrlObjErr = api.GetControllerInstance(ctrlName)
			if ctrlObjErr != nil {
				fmt.Println("get object error")
				web.Echo(w, []byte(ctrlObjErr.Error()))
				return
			}
		case "invalid":
			web.Echo(w, []byte("module name "+" is invalid."))
			return

		}
		method := web.GetMethodName(r)

		method = strings.Title(method)
		fmt.Println("method is " + method)
		f := reflect.ValueOf(ctrlObj).MethodByName(method)

		if !f.IsValid() {
			web.Echo(w, []byte(fmt.Sprintf("method `%s` is invalid.", method)))
			log.Logx.LogWarn("method `%s` is invalid.", method)
			return
		}
		function, _ := f.Interface().(func(http.ResponseWriter, *http.Request))
		function(w, r)

	})
	server := http.Server{Addr: ":" + port, Handler: mux}

	go func() {
		log.Logx.LogInfo("HttpServer is started, port %s", port)
		if err := server.ListenAndServe(); err != nil {
			log.Logx.LogError(err.Error())
		}
	}()

	sig := <-stopChan
	log.Logx.LogInfo("Receive signal: %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Logx.LogError("Server shutdown: %s", err.Error())
	} else {
		log.Logx.LogInfo("Server gracefully stopped.")
	}

}
