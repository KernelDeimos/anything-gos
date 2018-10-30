package interp_a

import (
	"errors"
	"fmt"
)

func BuiltinFormat(args []interface{}) ([]interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("format requires at least 1 argument")

	}
	args0, ok := args[0].(string)
	if !ok {
		return nil, errors.New("first argument must be string")
	}
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
