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

func (evaluator HybridEvaluator) RunEntry(
	entry HybridEvaluatorEntry,
	args []interface{},
) ([]interface{}, error) {
	switch entry.Tag {
	case EntryIsEvaluator:
		// Don't need to evaluate anything - pass all args literally
		return entry.Op(args)
	case EntryIsOperation:
		// Need to evaluate arguments before calling
		results := []interface{}{}
		for _, arg := range args {
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
		argsToPass = append(argsToPass, args...)

		// Call operation (the function)
		return entry.Op(argsToPass)

	case EntryIsNone:
		return args, errors.New("not found")
	}
	return nil, errors.New("evaluator: unrecognized operation tag: " +
		fmt.Sprint(entry.Tag))
}

type ErrorFunctionNotFound struct {
	FunctionName string
}

func (err ErrorFunctionNotFound) Error() string {
	return fmt.Sprintf("could not find function '%s'", err.FunctionName)
}

func (evaluator HybridEvaluator) GetEntry(args []interface{},
) (string, HybridEvaluatorEntry, bool, error) {
	nilEntry := HybridEvaluatorEntry{}
	if len(args) == 0 {
		return "", nilEntry, false, nil
	}
	first := args[0]
	switch opToRun := first.(type) {
	case string:
		if opToRun == "__debug_listmethods" {
			return opToRun, HybridEvaluatorEntry{
				Tag: EntryIsOperation,
				Op:  evaluator.OpListOperations,
			}, true, nil
		}
		entry, exists := evaluator.functionsMap[opToRun]
		if !exists {
			return opToRun, nilEntry, false, nil
		}
		return opToRun, entry, true, nil

	case Operation:
		return "__anonymous__", HybridEvaluatorEntry{
			Tag: EntryIsOperation,
			Op:  opToRun,
		}, true, nil

	default:
		return "", nilEntry, false, errors.New(
			"evaluator: only string types are supported")
	}
}

func (evaluator HybridEvaluator) evaluate(
	args []interface{},
) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}

	name, entry, exists, err := evaluator.GetEntry(args)
	if err != nil {
		return nil, err
	}
	if !exists {
		if evaluator.defaultBehaviour == nil {
			return nil, ErrorFunctionNotFound{FunctionName: name}
		} else {
			return evaluator.RunEntry(*(evaluator.defaultBehaviour), args)
		}
	}
	return evaluator.RunEntry(entry, args[1:])

}

func (evaluator HybridEvaluator) MakeChild() HybridEvaluator {
	child := HybridEvaluator{
		functionsMap:     map[string]HybridEvaluatorEntry{},
		defaultBehaviour: evaluator.defaultBehaviour,
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

func (evaluator HybridEvaluator) AddEvaluator(
	name string, function Operation,
) {
	evaluator.functionsMap[name+"-evaluate"] = HybridEvaluatorEntry{
		Tag: EntryIsOperation,
		Op:  function,
	}
	evaluator.functionsMap[name] = HybridEvaluatorEntry{
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

func (evaluator HybridEvaluator) OpListOperations(
	args []interface{},
) ([]interface{}, error) {
	entries := []interface{}{}
	for entry := range evaluator.functionsMap {
		entries = append(entries, interface{}(entry))
	}
	return entries, nil
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

func (evaluator HybridEvaluator) OpAddEvaluator(
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
	evaluator.functionsMap[name] = HybridEvaluatorEntry{
		Op:  function,
		Tag: EntryIsEvaluator,
	}
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
