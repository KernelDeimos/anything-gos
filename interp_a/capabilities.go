package interp_a

// TODO: (stretch goal)
// Interfaces like this could be generated from function definitions, which
// would be really cool. For instance, "CanOpEvaluateLikeHybridEvaluator"

type CanMakeChildEvaluator interface {
	MakeChild() CanEvaluate
}

type CanEvaluate interface {
	OpEvaluate(args []interface{}) ([]interface{}, error)
}
