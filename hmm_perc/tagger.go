package hmm_perc

import (
	"fmt"
	"os"
)

type Tagger struct {
	X     [][]string
	Nodes [][]Node

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
	if !IsTheSame(*y, pred) {
		this.Update(*y, pred)
	}
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
	//         fmt.Println(len(x)-1, "best", this.Output.Elems[n.Y], n.BestScore)
	i := len(this.X) - 2
	for {
		if i < 0 {
			break
		}
		n = n.Prev
		//                 fmt.Println(i, "best", this.Output.Elems[n.Y], n.BestScore)
		//                 pred[i] = this.LabelIndex.Elems[n.Y]
		pred[i] = this.feature_index.Output.Elems[n.Y]
		i--
	}
	return pred
}

func (this *Tagger) Update(y, pred []string) {
	var sign float64
	sign = 1.
	for i := 0; i < len(y); i++ {
		j := this.feature_index.Output.GetId(y[i])
		fs := this.Nodes[i][j].Fs
		for _, fid := range fs {
			this.feature_index.NodeWeight[j][fid] += 1.
		}

		lpath := this.Nodes[i][j].LPath
		for _, p := range lpath {
			y1 := p.LNode.Y
			y2 := p.RNode.Y
			if y2 == this.feature_index.Output.GetId(y[i]) &&
				y1 == this.feature_index.Output.GetId(y[i-1]) {
				for _, fid := range p.Fs {
					offset1 := y1*this.feature_index.Output.Size() + y2
					this.feature_index.EdgeWeight[offset1][fid] += 1.
				}
				break
			}
		}

		k := this.feature_index.Output.GetId(pred[i])
		fs = this.Nodes[i][k].Fs
		for _, fid := range fs {
			this.feature_index.NodeWeight[k][fid] -= 1
		}
		lpath = this.Nodes[i][k].LPath
		for _, p := range lpath {
			y1 := p.LNode.Y
			y2 := p.RNode.Y
			if y2 == this.feature_index.Output.GetId(pred[i]) &&
				y1 == this.feature_index.Output.GetId(pred[i-1]) {
				for _, fid := range p.Fs {
					offset1 := y1*this.feature_index.Output.Size() + y2
					this.feature_index.EdgeWeight[offset1][fid] -= sign
				}
				break
			}
		}
	}
}

func (this *Tagger) Save(model_file string) {
	this.feature_index.Save(model_file)
}

func a() {
	fmt.Println("AAA")
	os.Exit(1)
}