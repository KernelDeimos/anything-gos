package interp_a

//go:generate genfor-interp-a $GOFILE

import (
	"errors"
	"fmt"
)

type HybridEvaluatorEntryTag int

const (
	// EntryIsOperation indicates that the function expects its arguments to be
	// evaluated before the function is called.
	EntryIsOperation HybridEvaluatorEntryTag = iota

	// EntryIsEvaluator indicates that the function does not expect its
	// arguments to be evaluated, and does not need any evaluator function.
	EntryIsEvaluator

	// EntryIsControl indicates that the function does not expect its arguments
	// to be evaluated, but the function should always recieve the default
	// evaluator as the first parameter implicitly.
	EntryIsControl

	// EntryIsNone indicates that nothing exists in the entry
	// (currently only used for defaultBehaviour)
	EntryIsNone
)

type HybridEvaluatorEntry struct {
	Tag HybridEvaluatorEntryTag
	Op  Operation
}

type HybridEvaluator struct {
	functionsMap     map[string]HybridEvaluatorEntry
	defaultBehaviour *HybridEvaluatorEntry
}

func (evaluator HybridEvaluator) runEntry(
	entry HybridEvaluatorEntry,
	args []interface{},
) ([]interface{}, error) {
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
	case EntryIsControl:
		// Add default evaluator as first argument
		argsToPass := []interface{}{evaluator}

		// Add remaining arguments
		argsToPass = append(argsToPass, args[1:]...)

		// Call operation (the function)
		return entry.Op(argsToPass)
	}
	return nil, errors.New("evaluator: unrecognized operation tag")
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
			if evaluator.defaultBehaviour == nil {
				return nil, fmt.Errorf("could not find function '%s'", opToRun)
			} else {
				return evaluator.runEntry(*(evaluator.defaultBehaviour), args)
			}
		}
		return evaluator.runEntry(entry, args)

	default:
		return nil, errors.New("evaluator: only string types are supported")
	}
}

func (evaluator HybridEvaluator) MakeChild() HybridEvaluator {
	child := HybridEvaluator{
		functionsMap: map[string]HybridEvaluatorEntry{},
	}
	// TODO: This should eventually be optimised to use a linked list of
	//       hashmaps. This will require modification of this initializer
	//       and other constructors. get() and define() operations should be
	//       used to contain the logic for finding the most recent function
	//       definition for a specified namevalue down the chain.
	for key, f := range evaluator.functionsMap {
		child.functionsMap[key] = f
	}

	return child
}

func (evaluator HybridEvaluator) SetDefaultBehaviour(
	op Operation,
	tag HybridEvaluatorEntryTag,
) {
	evaluator.defaultBehaviour.Op = op
	evaluator.defaultBehaviour.Tag = tag
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

func (evaluator HybridEvaluator) GetOperation(
	name string,
) (Operation, HybridEvaluatorEntryTag, bool) {
	entry, found := evaluator.functionsMap[name]
	if !found {
		return nil, -1, false
	}
	return entry.Op, entry.Tag, true
}

func (evaluator HybridEvaluator) OpAddOperation(
	args []interface{},
) ([]interface{}, error) {
	//::gen verify-args add-operation name string function Operation
	if len(args) < 2 {
		return nil, errors.New("add-operation requires at least 2 arguments")
	}

	var name string
	var function Operation
	{
		var ok bool
		name, ok = args[0].(string)
		if !ok {
			return nil, errors.New("add-operation: argument 0: name; must be type string")
		}
		function, ok = args[1].(Operation)
		if !ok {
			return nil, errors.New("add-operation: argument 1: function; must be type Operation")
		}
	}
	//::end
	evaluator.AddOperation(name, function)
	return []interface{}{}, nil
}

func (evaluator HybridEvaluator) OpGetOperation(
	args []interface{},
) ([]interface{}, error) {
	//::gen verify-args add-operation name string
	if len(args) < 1 {
		return nil, errors.New("add-operation requires at least 1 arguments")
	}

	var name string
	{
		var ok bool
		name, ok = args[0].(string)
		if !ok {
			return nil, errors.New("add-operation: argument 0: name; must be type string")
		}
	}
	//::end
	op, _, found := evaluator.GetOperation(name)
	if !found {
		return nil, fmt.Errorf("could not find %s", name)
	}
	return []interface{}{op}, nil
}

func NewHybridEvaluator(
	functionsMap map[string]HybridEvaluatorEntry,
) (HybridEvaluator, error) {
	evaluator := HybridEvaluator{
		functionsMap: functionsMap,
		defaultBehaviour: &HybridEvaluatorEntry{
			Tag: EntryIsNone,
		},
	}

	return evaluator, nil
}
