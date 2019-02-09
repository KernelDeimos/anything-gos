package permutations

import (
	"errors"
)

type PermuSeq []bool

// Validate returns true of the length of the sequence is the result of an
// exponent of 2 (-1). This is a requirement for some of the methods to behave
// correctly.
func (ps PermuSeq) Validate() bool {
	if len(ps) <= 0 {
		return false
	}
	l := len(ps) + 1
	return (l & (l - 1)) == 0
}

func (ps PermuSeq) ToBinaryBoolTree() *BinaryBoolNode {
	var root *BinaryBoolNode
	type leafInfo struct {
		pointerPointer **BinaryBoolNode
		depth          int
	}
	nextLeaf := []leafInfo{}

	// Populate the tree with some BinaryBoolNodes
	for i := 0; i < len(ps); i++ {
		n := BinaryBoolNode{
			Value: ps[i],
			Left:  nil,
			Rite:  nil,
		}
		var nextDepth int
		if i == 0 {
			root = &n // yeah, I know; too lazy
			nextDepth = 1
		} else {
			leaf := nextLeaf[0]
			nextLeaf = nextLeaf[1:] // not a stack, press f for respect
			nextDepth = leaf.depth + 1
			*(leaf.pointerPointer) = &n
		}
		nextLeaf = append(nextLeaf, leafInfo{
			pointerPointer: &(n.Left),
			depth:          nextDepth,
		})
		nextLeaf = append(nextLeaf, leafInfo{
			pointerPointer: &(n.Rite),
			depth:          nextDepth,
		})
	}

	return root
}

func (ps PermuSeq) ToRelative() (PermuSeq, error) {
	if !ps.Validate() {
		return nil, errors.New("invalid")
	}

	root := ps.ToBinaryBoolTree()
	root.SwapTrues()
	return root.ToPermuSeq(), nil
}

func (ps PermuSeq) Permutate(input []interface{}) {
}
