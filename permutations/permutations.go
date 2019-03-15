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

func (ps PermuSeq) Permutate(input []interface{}) ([]interface{}, error) {
	if !ps.Validate() {
		return nil, errors.New("must be a valid permutation sequence")
	}

	// Calculate expected length of input (permutation key length + 1)
	len2 := len(ps) + 1

	// Validate or modify input length to match permutation sequence
	{
		if len2 < len(input) {
			return nil, errors.New("input slice too large")
		}
		for len(input) < len2 {
			input = append(input, nil)
		}
	}

	// Queue of things to swap
	// - each item: 2 indicies
	swapQueue := [][]int{}

	// Iteration count (used to read permutation sequence correctly)
	iter := 0

	// Initial condition: queue contains a task to swap the first two
	//                    halves of the input list.
	swapQueue = append(swapQueue, []int{0, len2})

	// Perform swapping algorithm
	for {
		// Stop if queue is empty
		if len(swapQueue) == 0 {
			break
		}

		// Pop the next swapping task
		swapTask := swapQueue[0]
		swapQueue = swapQueue[1:]

		// Calculate halfway point
		hwpoint := swapTask[0] + (swapTask[1]-swapTask[0])/2

		// Perform the swap task if the permutation sequence says to
		if ps[iter] { // ps[iter] is true if these halves should be swapped
			// Create slice for each half
			firstHalf := input[swapTask[0]:hwpoint]
			secndHalf := input[hwpoint:swapTask[1]]

			// Perform swap
			output := []interface{}{}
			output = append(output, input[0:swapTask[0]]...)
			output = append(output, secndHalf...)
			output = append(output, firstHalf...)
			output = append(output, input[swapTask[1]:]...)

			input = output
		}

		// If this task was to swap only 2 items, do not add tasks at this
		// iteration (could have also checked if hwpoint is an odd number)
		if swapTask[1]-swapTask[0] > 2 {

			// Add next tasks to queue
			swapQueue = append(swapQueue,
				[]int{swapTask[0], hwpoint},
				[]int{hwpoint, swapTask[1]},
			)
		}
		iter++
	}

	// At this point, the algorithm finished iterating; report no error
	return input, nil
}
