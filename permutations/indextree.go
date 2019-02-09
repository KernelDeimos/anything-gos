package permutations

// IndexTree is named this instead of IntegerTree just to remind myself that
// it's probably the only other tree I'll need.
type IndexTree struct {
	Value int
	Left  *BinaryBoolNode
	Rite  *BinaryBoolNode
}
