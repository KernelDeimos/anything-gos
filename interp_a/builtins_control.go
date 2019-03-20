package interp_a

import (
	"errors"
)

type ReturnState struct {
	Value []interface{}
	Ready bool
}

type ArgGetter struct {
	Value []interface{}
}

func MakeReturnState() *ReturnState {
	return &ReturnState{
		Value: []interface{}{},
		Ready: false,
	}
}

func (rs *ReturnState) Bind(eval HybridEvaluator) {
	eval.AddOperation("return",
		func(args []interface{}) ([]interface{}, error) {
			rs.Value = append(rs.Value, args...)
			rs.Ready = true
			return rs.Value, nil
		},
	)

	eval.AddOperation("return-ready",
		func(args []interface{}) ([]interface{}, error) {
			rs.Ready = true
			return rs.Value, nil
		},
	)

	eval.AddOperation("return-set",
		func(args []interface{}) ([]interface{}, error) {
			rs.Value = args
			return rs.Value, nil
		},
	)

	eval.AddOperation("return-push",
		func(args []interface{}) ([]interface{}, error) {
			rs.Value = append(rs.Value, args...)
			return []interface{}{rs.Ready}, nil
		},
	)

	eval.AddOperation("returning?",
		func(args []interface{}) ([]interface{}, error) {
			return []interface{}{rs.Ready}, nil
		},
	)

}

func MakeArgGetter() *ArgGetter {
	return &ArgGetter{
		Value: []interface{}{},
	}
}

func (rs *ArgGetter) Bind(eval HybridEvaluator) {
	eval.AddOperation("get",
		func(args []interface{}) ([]interface{}, error) {
			if len(args) == 1 {
				//::gen verify-args get_1 index int
				if len(args) < 1 {
					return nil, errors.New("get_1 requires at least 1 arguments")
				}

				var index int
				{
					var ok bool
					index, ok = args[0].(int)
					if !ok {
						return nil, errors.New("get_1: argument 0: index; must be type int")
					}
				}
				//::end

				if !(index < len(rs.Value)) {
					return nil, errors.New("index out of bounds")
				}
				value := rs.Value[index]
				return []interface{}{value}, nil
			}
			return rs.Value, nil
		},
	)
}

// BuiltinDo creates a sub-evaluator with return operations and evaluates each
// argument, discarding the return value of each argument.
//
// To concatenate return values, simply do:
// do (return-push (function-a)) (return-push (function-b)) (return)
func BuiltinDo(args []interface{}) ([]interface{}, error) {
	evalMaker := args[0].(HybridEvaluator)
	eval := evalMaker.MakeChild()

	rs := MakeReturnState()
	rs.Bind(eval)

	for _, arg := range args[1:] {
		switch value := arg.(type) {
		case string:
			// Ignore string (this is a comment)
		case []interface{}:
			// Skip empty list
			if len(value) < 1 {
				continue
			}
			// Evaluate list
			_, err := eval.OpEvaluate(value)
			// Check for error
			if err != nil {
				return rs.Value, err
			}
			// Check for return
			if rs.Ready {
				return rs.Value, nil
			}
		}
	}

	return rs.Value, nil
}

// BuiltinApply applies arguments to a section of code which wants to receive
// arguments
func BuiltinApply(args []interface{}) ([]interface{}, error) {
	if len(args) < 2 {
		return args, errors.New("Applying without arguments is nonsense")
	}
	if len(args) < 3 {
		return args, errors.New("Applying values to nothing is nonsense")
	}
	evalMaker := args[0].(HybridEvaluator)
	eval := evalMaker.MakeChild()

	ag := MakeArgGetter()
	input, ok := args[1].([]interface{})
	if !ok {
		return nil, errors.New("Input is not okay")
	}

	ag.Value = input
	ag.Bind(eval)

	applyTargetParameter := args[2:]
	applyTarget, err := evalMaker.OpEvaluate(applyTargetParameter)
	if err != nil {
		return applyTargetParameter, err
	}

	return eval.OpEvaluate(applyTarget)
}
