package domain

type NodeAI struct {
	isUsed bool
}

func NewNodeAI(APIKey string, isUsed bool) *NodeAI {
	return &NodeAI{
		isUsed: false,
	}
}
