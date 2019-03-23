package interp_a

import (
	"errors"
	"fmt"
	"strconv"
)

func BuiltinInt(args []interface{}) ([]interface{}, error) {
	result := []interface{}{}
	for _, arg := range args {
		iVal, err := strconv.ParseInt(fmt.Sprint(arg), 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, int(iVal))
	}
	return result, nil
}

func BuiltinLess(args []interface{}) ([]interface{}, error) {
	result := []interface{}{}
	if len(args) < 2 {
		return nil, errors.New("undefined output for < with one value")
	}
	for i := range args {
		if i == 0 {
			continue
		}
		leftInt, ok1 := args[i-1].(int)
		righInt, ok2 := args[i].(int)
		if !ok1 || !ok2 {
			return args, errors.New("<: non-integer parameters used")
		}

		if righInt < leftInt {
			return []interface{}{0}, nil
		}

		return []interface{}{1}, nil
	}
	return result, nil
}

func intReduce(args []interface{}, r func(values ...int) int,
) ([]interface{}, error) {
	integers := []int{}
	for i := range args {
		value, ok := args[i].(int)
		if !ok {
			return nil, errors.New("expected integer")
		}
		integers = append(integers, value)
	}
	return []interface{}{r(integers...)}, nil
}

func BuiltinAdd(args []interface{}) ([]interface{}, error) {
	return intReduce(args, func(values ...int) int {
		s := 0
		for _, v := range values {
			s += v
		}
		return s
	})
}

func BuiltinSubtract(args []interface{}) ([]interface{}, error) {
	return intReduce(args, func(values ...int) int {
		s := 0
		for _, v := range values {
			s -= v
		}
		return s
	})
}

func BuiltinMultiply(args []interface{}) ([]interface{}, error) {
	return intReduce(args, func(values ...int) int {
		s := 1
		for _, v := range values {
			s *= v
		}
		return s
	})
}

func BuiltinDivide(args []interface{}) ([]interface{}, error) {
	return intReduce(args, func(values ...int) int {
		s := 1
		for i, v := range values {
			if i%2 == 0 {
				s *= v
			} else {
				s /= v
			}
		}
		return s
	})
}
