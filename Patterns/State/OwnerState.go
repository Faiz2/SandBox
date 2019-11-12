package State

import (
	"SandBox/Model"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"strings"
)

type OwnerState struct {
	context *AuthContext
	db *BmMongodb.BmMongodb
	rd *BmRedis.BmRedis
	fileModel Model.FileMetaDatum
}

func (os *OwnerState) DoExecute() (interface{}, error) {
	if len(os.fileModel.Mod) == 0 { os.ChangeState(false); return os.context.state.DoExecute() }

	mod := os.fileModel.Mod[1:4]
	os.ChangeState(strings.Contains(mod, os.context.action))
	return os.context.state.DoExecute()
}

func (os *OwnerState) ChangeState(status bool) {
	if status { os.context.state = &EndState{}
	} else { os.context.state = &NonePermissionsState{} }
}