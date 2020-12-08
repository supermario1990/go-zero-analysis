package main

import (
	"flag"

	"shorturl/api/internal/config"
	"shorturl/api/internal/handler"
	"shorturl/api/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

// 配置文件
// 支持三种后缀：json,yaml,yml
var configFile = flag.String("f", "etc/shorturl-api.yaml", "the config file")

func main() {
	flag.Parse()

	// 加载配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// ctx 由配置文件、zrpc客户端组成
	ctx := svc.NewServiceContext(c)

	// 创建rest api 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册路由
	handler.RegisterHandlers(server, ctx)

	// 启动server
	server.Start()
}
