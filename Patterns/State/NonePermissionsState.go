package State

import "errors"

type NonePermissionsState struct {}

func (nps *NonePermissionsState) DoExecute() (interface{}, error) {
	return nil, errors.New("None Permissions")
}
func (nps *NonePermissionsState) ChangeState(status bool) {
	panic("This is the end point")
}