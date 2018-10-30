package interp_a

import (
	"errors"
	"fmt"
)

type Operation func(args []interface{}) ([]interface{}, error)

func ListToMap(items []interface{}) (map[string]interface{}, error) {
	if len(items)%2 != 0 {
		return nil, errors.New("list-to-map needs even argument count")
	}

	mapToReturn := map[string]interface{}{}

	for i := 0; i < len(items); i += 2 {
		key, keyIsString := items[i].(string)
		value := items[i+1]

		if !keyIsString {
			return nil, errors.New("only string keys are supported currently")
		}

		mapToReturn[key] = value
	}

	return mapToReturn, nil
}

func NewOperationMux(
	functionsMap map[string]Operation,
) (Operation, error) {
	var op Operation

	op = func(args []interface{}) ([]interface{}, error) {
		if len(args) == 0 {
			return nil, nil
		}
		first := args[0]
		switch opToRun := first.(type) {
		case string:
			opFunc, exists := functionsMap[opToRun]
			if !exists {
				return nil, fmt.Errorf("could not find function '%s'", opToRun)
			}
			return opFunc(args[1:])

		}
		return nil, nil
	}

	return op, nil
}
