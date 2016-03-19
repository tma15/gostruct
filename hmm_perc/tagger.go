package hmm_perc

import (
	"container/heap"
	"fmt"
	"os"
)

type Tagger struct {
	X     [][]string
	Nodes [][]Node
	EOS   Node

	feature_index FeatureIndex
}

func NewTagger() Tagger {
	this := Tagger{
		Nodes: make([][]Node, 0, 10),
	}
	return this
}

/* readして観測列xを作成 */
func (this *Tagger) SetX(x [][]string) {
	this.X = x
}

func (this *Tagger) SetFeatureIndex(fi FeatureIndex) {
	this.feature_index = fi
}

func IsTheSame(y1, y2 []string) bool {
	n := len(y1)
	for i := 0; i < n; i++ {
		if y1[i] != y2[i] {
			return false
		}
	}
	return true
}

func (this *Tagger) Fit(x *[][]string, y *[]string) {
	this.SetX(*x)
	this.feature_index.BuildFeatures(this)
	this.Viterbi()
	pred := this.BackTrack()
	this.Update(*y, pred)
}

func (this *Tagger) Predict(x [][]string) []string {
	this.SetX(x)
	this.feature_index.BuildFeatures(this)
	this.Viterbi()
	pred := this.BackTrack()
	return pred
}

func (this *Tagger) Viterbi() {
	for i := 0; i < len(this.X); i++ {
		for j, _ := range this.Nodes[i] {
			node := this.Nodes[i][j]
			lpath := node.LPath
			bestscore := -1e+10
			/* bestnodeは現在のnodeに繋がる左のノードで最もscoreが良いノード */
			var bestnode Node
			for _, p := range lpath {
				if i == 1 {
					p.LNode.BestScore = p.LNode.Score
				}
				s := p.LNode.BestScore + p.Score + node.Score
				if bestscore < s {
					bestnode = *p.LNode
					bestscore = s
				}
			}
			this.Nodes[i][j].Prev = &bestnode
			this.Nodes[i][j].BestScore = bestscore
		}
	}
}

func (this *Tagger) BackTrack() []string {
	last := len(this.X) - 1
	var bestnode *Node
	bestscore := -1e+10
	for j, _ := range this.Nodes[last] {
		node := this.Nodes[last][j]
		if node.Score > bestscore {
			bestscore = node.Score
			bestnode = &node
		}
	}
	if bestnode == nil {
		panic("best node is nil")
	}

	pred := make([]string, len(this.X), len(this.X))
	var n *Node = bestnode
	pred[len(this.X)-1] = this.feature_index.Output.Elems[n.Y]
	fmt.Println(len(this.X)-1, "best", this.feature_index.Output.Elems[n.Y],
		n.BestScore, n.Prev)
	i := len(this.X) - 2
	for {
		if i < 0 {
			break
		}
		n = n.Prev
		fmt.Println(i, "best", this.feature_index.Output.Elems[n.Y], n.BestScore, n.Prev)
		//                 pred[i] = this.LabelIndex.Elems[n.Y]
		pred[i] = this.feature_index.Output.Elems[n.Y]
		i--
	}
	return pred
}

func (this *Tagger) Update(y, pred []string) {
	if !IsTheSame(y, pred) {
		for i := 0; i < len(y); i++ {
			j := this.feature_index.Output.GetId(y[i])
			k := this.feature_index.Output.GetId(pred[i])
			/* true */
			fs := this.Nodes[i][j].Fs
			for _, fid := range fs {
				this.feature_index.NodeWeight[j][fid] += 1.
			}
			/* predict */
			fs = this.Nodes[i][k].Fs
			for _, fid := range fs {
				this.feature_index.NodeWeight[k][fid] -= 1.
			}

			if i > 0 {
				p2 := this.feature_index.Output.GetId(pred[i])
				p1 := this.feature_index.Output.GetId(pred[i-1])
				t2 := this.feature_index.Output.GetId(y[i])
				t1 := this.feature_index.Output.GetId(y[i-1])
				this.feature_index.EdgeWeight[t1][t2] += 1.
				this.feature_index.EdgeWeight[p1][p2] -= 1.
			}

		}
	}

}

func (this *Tagger) BackwardAstar(N int, nodes [][]*Node, eos *Node) [][]int {
	pqueue := make(PriorityQueue, 0, 10)
	heap.Init(&pqueue)

	eos.PathTotalScore = 0.
	heap.Push(&pqueue, eos)

	var result []*Node = make([]*Node, N, N)
	var n int = 0

	for {
		if pqueue.IsEmpty() {
			break
		}
		var node *Node = pqueue.Pop().(*Node)
		//         fmt.Println("POP", node.X, node.Y)

		if node.X == 0 { // is bos
			result[n] = node
			n += 1
		} else {
			for _, p := range node.LPath {
				prev := *p.LNode // copy
				prev.GoalScore = node.GoalScore + p.Score
				prev.PathTotalScore = prev.BestScore + prev.GoalScore
				prev.Next = node
				//                 fmt.Println("PUSH", prev.X, prev.Y, prev.Score)
				heap.Push(&pqueue, &prev)
			}
		}

		if n >= N {
			break
		}
	}

	var results [][]int = make([][]int, 0, N)
	for i, n := range result {
		var result []int = make([]int, 0, len(nodes))
		for {
			n = n.Next
			if n.X == eos.X {
				break
			}

			fmt.Println(fmt.Sprintf("%d-best: %d", i+1, n.Y))
			result = append(result, n.Y)
		}
		fmt.Println("")
	}
	return results
}

func (this *Tagger) Save(model_file string) {
	this.feature_index.Save(model_file)
}

func a() {
	fmt.Println("AAA")
	os.Exit(1)
}
