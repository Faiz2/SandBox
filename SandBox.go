// Package Main Program Entrance
package main

import (
	"SandBox/Factory"
	"SandBox/env"
	"fmt"
	"github.com/PharbersDeveloper/bp-go-lib/log"
	"github.com/alfredyang1986/BmServiceDef/BmApiResolver"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/alfredyang1986/BmServiceDef/BmPodsDefine"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"github.com/rs/cors"
	"net/http"
	"os"
)

func main() {
	// 本地调试打开
	env.SetLocalEnv()

	// 设置启动参数
	env.SetStartingParameter()

	version := "v0"
	prodEnv := "SANDBOX_HOME"
	log.NewLogicLoggerBuilder().Build().Info("SandBoxPods begins, version = ", version)

	fac := Factory.Table{}
	pod := BmPodsDefine.Pod{Name: "new SandBox", Factory: fac}
	home := os.Getenv(prodEnv)
	pod.RegisterSerFromYAML(home + "/resource/service-def.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig(prodEnv)

	addr := fmt.Sprint(bmRouter.Host, ":", bmRouter.Port)
	log.NewLogicLoggerBuilder().Build().Info("Listening on ", addr)

	api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)
	pod.RegisterAllFunctions(version, api)
	pod.RegisterAllMiddleware(api)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	handler := api.Handler().(*httprouter.Router)
	pod.RegisterPanicHandler(handler)
	err := http.ListenAndServe(":" + bmRouter.Port, c.Handler(handler))
	log.NewLogicLoggerBuilder().Build().Error(err)
}