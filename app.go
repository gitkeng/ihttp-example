package main

import (
	"github.com/gitkeng/ihttp"
	"github.com/gitkeng/ihttp/log"
)

func main() {
	cfgLocation := "./conf/app_setting.yaml"
	conf, err := ihttp.NewConfig(cfgLocation)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	logCfg, isLogCfgFound := conf.GetLogConfig()
	if !isLogCfgFound {
		log.Fatalf("Error: log config not found")
	}
	apiCfg, isAPICfgFound := conf.GetAPIConfig()
	if !isAPICfgFound {
		log.Fatalf("Error: log apit not found")
	}
	dbCfgs, isDBCfgFound := conf.GetDBConfigs()
	if !isDBCfgFound {
		log.Fatalf("Error: db config not found")
	}
	redisCfgs, isRedisCfgFound := conf.GetRedisConfigs()
	if !isRedisCfgFound {
		log.Fatalf("Error: redis config not found")
	}
	
	ms, err := ihttp.New(
		ihttp.WithAPIConfig(apiCfg),
		ihttp.WithLogConfig(logCfg),
		ihttp.WithDBConfigs(dbCfgs...),
		ihttp.WithRedisConfigs(redisCfgs...),
	)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	// register handler
	registerEmployeeHandler(ms)

	if err := ms.Start(); err != nil {
		log.Warnf("Error: %s", err.Error())
	}
}
