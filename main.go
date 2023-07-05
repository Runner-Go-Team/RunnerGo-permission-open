package main

import (
	"flag"
	"fmt"
	"permission-open/internal"
	"permission-open/internal/app/router"
	"permission-open/internal/pkg/conf"

	"github.com/gin-gonic/gin"
)

var readConfMode int
var configFile string

func main() {
	flag.IntVar(&readConfMode, "m", 0, "读取环境变量还是读取配置文件")
	flag.StringVar(&configFile, "c", "./configs/dev.yaml", "app config file.")
	flag.Parse()

	internal.InitProjects(readConfMode, configFile)

	//pyroscope.Start(
	//	pyroscope.Config{
	//		ApplicationName: "RunnerGo-permission-open",
	//		ServerAddress:   "http://192.168.1.205:4040/",
	//		Logger:          pyroscope.StandardLogger,
	//		ProfileTypes: []pyroscope.ProfileType{
	//			pyroscope.ProfileCPU,
	//			pyroscope.ProfileAllocObjects,
	//			pyroscope.ProfileAllocSpace,
	//			pyroscope.ProfileInuseObjects,
	//			pyroscope.ProfileInuseSpace,
	//		},
	//	})

	r := gin.New()
	router.RegisterRouter(r)

	// 全局参数企业相关信息
	internal.InitSomething()

	if err := r.Run(fmt.Sprintf(":%d", conf.Conf.Http.Port)); err != nil {
		panic(err)
	}
}
