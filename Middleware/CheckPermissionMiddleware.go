package Middleware

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type CheckPermissionMiddleware struct {
	Args []string
	db   *BmMongodb.BmMongodb
	rd   *BmRedis.BmRedis
}

func (cpm CheckPermissionMiddleware) NewCheckPermissionMiddleware(args ...interface{}) CheckPermissionMiddleware {
	var r *BmRedis.BmRedis
	var m *BmMongodb.BmMongodb
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				}
			}
		} else if i == 1 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}
	return CheckPermissionMiddleware{
		Args: ag,
		db:   m,
		rd:   r,
	}
}

func (cpm CheckPermissionMiddleware) DoMiddleware(c api2go.APIContexter, w http.ResponseWriter, r *http.Request) {
	//if len(r.Header.Get("Target")) > 0 {
	//	log.NewLogicLoggerBuilder().Build().Info(r.Header.Get("Action"))
	//	log.NewLogicLoggerBuilder().Build().Info(r.Header.Get("Who"))
	//	log.NewLogicLoggerBuilder().Build().Info(r.Header.Get("Target"))
	//	who := map[string]string {}
	//	context := State.AuthContext{}
	//	err :=	json.Unmarshal([]byte(r.Header.Get("Who")), &who)
	//	if err != nil { panic(err.Error()) }
	//	context.NewAuthContext(r.Header.Get("Target"), r.Header.Get("Action"), who,  cpm.db, cpm.rd)
	//	_, err = context.DoExecute()
	//	if err != nil { panic(err.Error()) }
	//}
	//log.NewLogicLoggerBuilder().Build().Info("Permission Middleware")
}
