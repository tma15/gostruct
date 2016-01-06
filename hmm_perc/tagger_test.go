package hmm_perc

import (
	"fmt"
	"testing"
)

func TestViterbi(t *testing.T) {
	x := [][]string{
		[]string{"the", "DT"},
		[]string{"current", "JJ"},
	}
	y := []string{"B-NP", "I-NP"}

	graph := `
	a       noun
I-NP  node1 --- node3
	    \ /
	     x
	    / \
B-NP  node2 --- node4
	`
	fmt.Println(graph)
	node1 := NewNode()
	node1.X = 0 /*a*/
	node1.Y = 0 /*B-NP*/
	node1.Fs = []int{0}
	node1.Obs = x[0][0]
	node1.Score = 1

	node2 := NewNode()
	node2.X = 0 /*a*/
	node2.Y = 1 /*I-NP*/
	node2.Fs = []int{0}
	node2.Obs = x[1][0]
	node2.Score = 1

	node3 := NewNode()
	node3.X = 1 /*noun*/
	node3.Y = 0 /*B-NP*/
	node3.Fs = []int{1}
	node3.Obs = x[0][0]
	node3.Score = 0

	node4 := NewNode()
	node4.X = 1 /*noun*/
	node4.Y = 1 /*I-NP*/
	node4.Fs = []int{1}
	node4.Obs = x[1][0]
	node4.Score = 1

	path1 := Path{
		LNode: &node1,
		RNode: &node3,
	}
	path1.Add(&node1, &node3)

	path2 := Path{
		LNode: &node2,
		RNode: &node3,
	}
	path2.Add(&node2, &node3)

	path3 := Path{
		LNode: &node1,
		RNode: &node4,
	}
	path3.Add(&node1, &node4)

	path4 := Path{
		LNode: &node2,
		RNode: &node4,
	}
	path4.Add(&node2, &node4)

	nodes := [][]Node{
		[]Node{
			node1,
			node2,
		},
		[]Node{
			node3,
			node4,
		},
	}

	tagger := NewTagger()
	tagger.SetX(x)
	tagger.Nodes = nodes
	fi := NewFeatureIndex()
	fi.Open("./sample/example.tmp", "./sample/input.txt")
	tagger.SetFeatureIndex(fi)
	tagger.Viterbi()
	pred := tagger.BackTrack()
	fmt.Println(y, pred)
}
