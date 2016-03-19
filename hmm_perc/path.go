package hmm_perc

type Path struct {
	RNode *Node
	LNode *Node

	Fs []int

	Score float64
}

func (this *Path) Add(lnode, rnode *Node) {
	this.LNode = lnode
	rnode.LPath = append(rnode.LPath, this)
}
