package interp_a

type InterpreterFactoryA struct{}

func (ifa InterpreterFactoryA) MakeExec() Operation {
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
	//::end

	exe, err := NewHybridEvaluator(fmap)
	if err != nil {
		// This error only occurs if initialization code above is invalid
		panic(err)
	}

	exe.AddOperation("@", exe.OpAddOperation)

	return exe.OpEvaluate
}
