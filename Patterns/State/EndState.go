package State

type EndState struct {}

func (es *EndState) DoExecute() (interface{}, error) {
	return true, nil
}

func (es *EndState) ChangeState(status bool) {
	panic("This is the end point")
}