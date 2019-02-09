package interp_a

type ReturnState struct {
	Value []interface{}
	Ready bool
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
