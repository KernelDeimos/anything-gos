package permutations

// BinaryIntNode is named this instead of IntegerTree just to remind myself that
// it's probably the only other tree I'll need.
type BinaryIntNode struct {
	Value int
	Left  *BinaryBoolNode
	Rite  *BinaryBoolNode
}
