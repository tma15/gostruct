package hmm_perc

type Path struct {
	RNode *Node
	LNode *Node
	Fs    []int
	Score float64
}

func (this *Path) Add(lnode, rnode *Node) {
	//         lnode.RPath = append(lnode.RPath, this)
	rnode.LPath = append(rnode.LPath, this)
}
