# Aquarius

## Aquarius 是一个轻量级后台开发框架

### 启动方式


启动方式     ./Aquarius -server=HTTPServer  -port=8765

默认启动 HTTPServer 端口8765    //TCPserver正在实现中

### pb协议
```
message apiReq{
required string action = 1;
required string method = 2;
required string systemId = 3;
}

message apiResp{
required int32 code = 1;
required string errmsg = 2;
optional bytes data = 3;
}
```

### 主流程

获取httpserver对象->监听系统signal->设置pprof监控系统信息->进入统一服务入口/->执行处理流程

生成httpcontext->获取对应controller及执行动作Action对象

### 注册controller中action对象

添加新的controller时需在路由表router中注册，否则无法执行，其格式为map[classtype][classname]
在action对象的init方法中调用registerControllerName即可。

### 日志系统

分为 按大小切分，时间切分（天，小时，分钟）
```
LOG_SHIFT_BY_SIZE = iota
LOG_SHIFT_BY_DAY
LOG_SHIFT_BY_HOUR
LOG_SHIFT_BY_MINUTE

type Log struct {
	log_path       string
	log_prefix     string
	log_level      int
	log_num        int
	log_size       int64
	log_shift_type int
	file           *os.File
	log_filename   string
	mutex          *sync.Mutex
}
```

### 缓存系统

设置全局缓存，每30秒清理过期缓存 格式为map[string]cacheNode
```
type cacheNode struct {
	Data       interface{}
	Expiration time.Time
}
```

### config
必须放置在 path文件中指定的目录位置，以map形式存储在缓存中，可以在线更新。
### api服务流程

api统一服务入口->常规检查->鉴权->获取api处理器->返回封装信息
