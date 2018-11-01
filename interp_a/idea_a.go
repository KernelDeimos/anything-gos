package interp_a

type InterpreterFactoryA struct{}

// MakeEmpty makes an evaluator with only mutator functions
func (ifa InterpreterFactoryA) MakeEmpty() HybridEvaluator {
	exe, err := NewHybridEvaluator(map[string]HybridEvaluatorEntry{})
	if err != nil {
		// This error only occurs if initialization code above is invalid
		panic(err)
	}

	// Bind evaluator mutators
	exe.AddOperation(":", exe.OpAddOperation)
	exe.AddOperation("$", exe.OpGetOperation)

	return exe
}

// MakeExec makes the root operation for the interpreter, which is a
// HybridEvaluator with builtin functions already added to it.
func (ifa InterpreterFactoryA) MakeExec() HybridEvaluator {
	fmap := map[string]HybridEvaluatorEntry{}

	// Function to add a new operation (o)
	o := func(name string, function Operation) {
		fmap[name] = HybridEvaluatorEntry{
			Tag: EntryIsOperation,
			Op:  function,
		}
		fmap[name+"-literal"] = HybridEvaluatorEntry{
			Tag: EntryIsEvaluator,
			Op:  function,
		}
	}

	// Function to add a new evaluator (e)
	/*
		e := func(name string, function Operation) {
			fmap[name] = HybridEvaluatorEntry{
				Tag: EntryIsEvaluator,
				Op:  function,
			}
			fmap[name+"-evaluate"] = HybridEvaluatorEntry{
				Tag: EntryIsOperation,
				Op:  function,
			}
		}
	*/

	//::gen register-all-functions
	o("format", BuiltinFormat)
	o("cat", BuiltinCat)
	o("store", BuiltinStore)
	//::end

	exe, err := NewHybridEvaluator(fmap)
	if err != nil {
		// This error only occurs if initialization code above is invalid
		panic(err)
	}

	// Bind evaluator mutators
	exe.AddOperation(":", exe.OpAddOperation)
	exe.AddOperation("$", exe.OpGetOperation)

	return exe
}
