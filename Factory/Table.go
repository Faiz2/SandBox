package Factory

import (
	"SandBox/DataStorage"
	"SandBox/Handler"
	"SandBox/Middleware"
	"SandBox/Model"
	"SandBox/Resource"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
)

type Table struct{}

var MODEL_FACTORY = map[string]interface{}{
	"SandBoxIndex": 	Model.SandBoxIndex{},
	"FileMetaDatum": 	Model.FileMetaDatum{},
	"FileVersion": 		Model.FileVersion{},
}

var STORAGE_FACTORY = map[string]interface{}{
	"SandBoxIndexStorage": 		DataStorage.SandBoxIndexStorage{},
	"FileMetaDatumStorage": 	DataStorage.FileMetaDatumStorage{},
	"FileVersionStorage":  		DataStorage.FileVersionStorage{},
}

var RESOURCE_FACTORY = map[string]interface{}{
	"SandBoxIndexResource": 	Resource.SandBoxIndexResource{},
	"FileMetaDatumResource": 	Resource.FileMetaDatumResource{},
	"FileVersionResource": 		Resource.FileVersionResource{},
}

var FUNCTION_FACTORY = map[string]interface{}{
	"CommonPanicHandle": 		Handler.CommonPanicHandle{},
	"GetAccountDetailHandler": 	Handler.GetAccountDetailHandler{},
	"PutHDFSHandler": 	Handler.PutHDFSHandler{},
	"StreamOss2HDFSHandler": 	Handler.PutHDFSHandler{},
	"UpdateJobIDWithTraceIDHandler": 	Handler.PutHDFSHandler{},
	"Stream2HDFSFinishHandler": 	Handler.PutHDFSHandler{},
}
var MIDDLEWARE_FACTORY = map[string]interface{}{
	"CheckTokenMiddleware": Middleware.CheckTokenMiddleware{},
	"CheckPermissionMiddleware": Middleware.CheckPermissionMiddleware{},
}

var DAEMON_FACTORY = map[string]interface{}{
	"BmMongodbDaemon": BmMongodb.BmMongodb{},
	"BmRedisDaemon":   BmRedis.BmRedis{},
}

func (t Table) GetModelByName(name string) interface{} {
	return MODEL_FACTORY[name]
}

func (t Table) GetResourceByName(name string) interface{} {
	return RESOURCE_FACTORY[name]
}

func (t Table) GetStorageByName(name string) interface{} {
	return STORAGE_FACTORY[name]
}

func (t Table) GetDaemonByName(name string) interface{} {
	return DAEMON_FACTORY[name]
}

func (t Table) GetFunctionByName(name string) interface{} {
	return FUNCTION_FACTORY[name]
}

func (t Table) GetMiddlewareByName(name string) interface{} {
	return MIDDLEWARE_FACTORY[name]
}
