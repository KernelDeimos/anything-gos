package interp_a

//go:generate genfor-interp-a $GOFILE

import (
	"errors"
	"fmt"
)

func BuiltinFormat(args []interface{}) ([]interface{}, error) {
	//::gen verify-args format args0 string
	if len(args) < 1 {
		return nil, errors.New("format requires at least 1 arguments")
	}

	var args0 string
	{
		var ok bool
		args0, ok = args[0].(string)
		if !ok {
			return nil, errors.New("format: argument 0: args0; must be type string")
		}
	}
	//::end
	result := fmt.Sprintf(args0, args[1:]...)

	return []interface{}{result}, nil
}

func BuiltinCat(args []interface{}) ([]interface{}, error) {
	result := ""
	for _, arg := range args {
		result = result + fmt.Sprint(arg)
	}
	return []interface{}{result}, nil
}

func BuiltinStore(args []interface{}) ([]interface{}, error) {
	op := func(_ []interface{}) ([]interface{}, error) {
		return args, nil
	}
	return []interface{}{Operation(op)}, nil
}
