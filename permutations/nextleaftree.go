package permutations

type DepthTrackingBinaryBoolNode struct {
	Value    bool
	Depth    int
	Contents interface{}
	Left     *DepthTrackingBinaryBoolNode
	Rite     *DepthTrackingBinaryBoolNode
}

type DepthTrackingBinaryBoolTree struct {
	root      *DepthTrackingBinaryBoolNode
	nextNodes []*DepthTrackingBinaryBoolNode
}
