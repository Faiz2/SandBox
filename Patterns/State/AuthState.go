package State

type State interface {
	DoExecute() (interface{}, error)
	ChangeState(status bool)
}