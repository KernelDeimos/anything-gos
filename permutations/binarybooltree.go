package permutations

type BinaryBoolNode struct {
	Value bool
	Depth int
	Left  *BinaryBoolNode
	Rite  *BinaryBoolNode
}

func (tree *BinaryBoolNode) ToPermuSeq() []bool {
	seq := PermuSeq{}

	nextNode := []*BinaryBoolNode{}
	nextNode = append(nextNode, tree)
	for len(nextNode) > 0 {
		node := nextNode[0]
		nextNode = nextNode[1:]
		seq = append(seq, node.Value)

		if node.Left != nil {
			nextNode = append(nextNode, node.Left)
		}
		if node.Rite != nil {
			nextNode = append(nextNode, node.Rite)
		}
	}

	return seq
}

func (tree *BinaryBoolNode) SwapTrues() {
	// TODO: non-recursive implementation
	tree.SwapTruesRecursive()
}

func (tree *BinaryBoolNode) SwapTruesRecursive() {
	if tree == nil {
		return
	}
	if tree.Value {
		tmp := tree.Left
		tree.Left = tree.Rite
		tree.Rite = tmp
	}
	// Does not matter what order these are called in
	tree.Left.SwapTruesRecursive()
	tree.Rite.SwapTruesRecursive()
}
