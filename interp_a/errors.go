package interp_a

type ErrorType struct {
	inputArgs []interface{}
	nextArgs  []interface{}
	info      string
	nextError error
}

func (et ErrorType) Error() string {
	return et.info + ": " + et.nextError.Error()
}

func (et ErrorType) String() string {
	return "trace(" + et.info + "): " + et.nextError.Error()
}

func (et ErrorType) Return() ([]interface{}, error) {
	return []interface{}{
		et,
		et.String(),
		et.inputArgs,
		"->",
		et.nextArgs,
	}, et
}

func resultForError(info string, args, nextArgs []interface{}, ein error,
) ([]interface{}, error) {
	et := ErrorType{
		inputArgs: args,
		nextArgs:  nextArgs,
		info:      info,
		nextError: ein,
	}
	return et.Return()
}
