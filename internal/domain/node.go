package domain

type NodeAI struct {
	APIkey string
	isUsed bool
}

func NewNodeAI(APIKey string, isUsed bool) *NodeAI {
	return &NodeAI{
		APIkey: APIKey,
		isUsed: false,
	}
}
