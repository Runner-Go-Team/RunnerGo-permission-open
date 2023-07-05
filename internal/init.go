package internal

import (
	"context"
	"permission-open/internal/pkg/biz/log"
	"permission-open/internal/pkg/biz/proof"
	"permission-open/internal/pkg/conf"
	"permission-open/internal/pkg/dal"
	"permission-open/internal/pkg/event"
	"permission-open/internal/pkg/logic/auth"

	"go.uber.org/zap"
)

func InitProjects(readConfMode int, configFile string) {
	if readConfMode == 1 {
		conf.MustInitConfByEnv()
	} else {
		conf.MustInitConf(configFile)
	}

	//conf.MustInitConf()
	// 初始化各种中间件
	dal.MustInitMySQL()
	dal.MustInitMongo()
	//dal.MustInitElasticSearch()
	proof.MustInitProof()
	//dal.MustInitGRPC()
	dal.MustInitRedis()
	dal.MustInitBigCache()
	// 初始化logger
	zap.S().Debug("初始化logger")
	log.InitLogger()
}

func InitSomething() {
	_, _ = auth.CompanyRegister(context.Background(), conf.Conf.CompanyInitConfig.Account, conf.Conf.CompanyInitConfig.Password)

	// 注册 event
	_ = event.Register()
}
