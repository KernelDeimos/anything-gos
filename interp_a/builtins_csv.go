package interp_a

import (
	"errors"
	"strconv"
	"strings"
)

// Status: horrible, probably don't use

func BuiltinListToCsvletsN(args []interface{}) ([]interface{}, error) {
	strargs := []string{}

	//modified ::gen verify-args list-to-csvlets args0 float64
	if len(args) < 1 {
		return nil, errors.New("list-to-csvlets requires at least 1 arguments")
	}

	// TODO: need a much better way to handle this
	var args0i int64
	{
		var err error
		args0s, ok := args[0].(string)
		if !ok {
			return nil, errors.New("must be string")
		}
		args0i, err = strconv.ParseInt(args0s, 10, 64)
		if err != nil {
			return nil, errors.New("list-to-csvlets: argument 0: args0; must be type float64")
		}
	}
	//modified ::end

	for _, arg := range args[1:] {
		if strval, ok := arg.(string); ok {
			strargs = append(strargs, strval)
		} else {
			return nil, errors.New("all arguments must be strings")
		}
	}

	pairs := [][]string{}
	var currPair *[]string

	for i := 0; i < len(strargs); i++ {
		if 0 == (i % int(args0i)) {
			if currPair != nil {
				pairs = append(pairs, *currPair)
			}
			nextPair := []string{}
			currPair = &nextPair
		}
		*currPair = append(*currPair, strargs[i])
	}
	if currPair != nil {
		pairs = append(pairs, *currPair)
	}

	csvlets := []interface{}{}
	for _, pair := range pairs {
		csvlets = append(csvlets, strings.Join(pair, ","))
	}

	return csvlets, nil
}
