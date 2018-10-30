package interp_a

import (
	"errors"
	"fmt"
)

type HybridEvaluatorEntryTag int

const (
	EntryIsOperation HybridEvaluatorEntryTag = iota
	EntryIsEvaluator
)

type HybridEvaluatorEntry struct {
	Tag HybridEvaluatorEntryTag
	Op  Operation
}

type HybridEvaluator struct {
	functionsMap map[string]HybridEvaluatorEntry
}

func (evaluator HybridEvaluator) evaluate(
	args []interface{},
) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	first := args[0]
	switch opToRun := first.(type) {
	case string:
		entry, exists := evaluator.functionsMap[opToRun]
		if !exists {
			return nil, fmt.Errorf("could not find function '%s'", opToRun)
		}
		switch entry.Tag {
		case EntryIsEvaluator:
			// Don't need to evaluate anything - pass all args literally
			return entry.Op(args[1:])
		case EntryIsOperation:
			// Need to evaluate arguments before calling
			results := []interface{}{}
			for _, arg := range args[1:] {
				switch argToEvaluate := arg.(type) {
				case []interface{}:
					theseResults, err := evaluator.evaluate(
						argToEvaluate,
					)
					if err != nil {
						return nil, err
					}
					for _, item := range theseResults {
						results = append(results, item)
					}
				default:
					results = append(results, arg)
				}
			}
			return entry.Op(results)
		}
		return nil, errors.New("evaluator: unrecognized operation tag")
	default:
		return nil, errors.New("evaluator: only string types are supported")
	}
	return nil, nil
}

func (evaluator HybridEvaluator) OpEvaluate(
	args []interface{},
) ([]interface{}, error) {
	return evaluator.evaluate(args)
}

func (evaluator HybridEvaluator) AddOperation(
	name string, function Operation,
) {
	evaluator.functionsMap[name] = HybridEvaluatorEntry{
		Tag: EntryIsOperation,
		Op:  function,
	}
	evaluator.functionsMap[name+"-literal"] = HybridEvaluatorEntry{
		Tag: EntryIsEvaluator,
		Op:  function,
	}
}

func (evaluator HybridEvaluator) OpAddOperation(
	args []interface{},
) ([]interface{}, error) {
	//::gen verify-args add-operation name string function Operation
	if len(args) < 2 {
		return nil, errors.New("add-operation takes 2 arguments")
	}
	var name string
	var function Operation
	{
		var ok bool
		name, ok = args[0].(string)
		if !ok {
			return nil, errors.New("<name> must be a string")
		}
		function, ok = args[1].(Operation)
		if !ok {
			return nil, errors.New("<function> must be an operation")
		}
	}
	//::end
	evaluator.AddOperation(name, function)
	return []interface{}{}, nil
}

func NewHybridEvaluator(
	functionsMap map[string]HybridEvaluatorEntry,
) (HybridEvaluator, error) {
	evaluator := HybridEvaluator{
		functionsMap: functionsMap,
	}

	return evaluator, nil
}
