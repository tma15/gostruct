package hmm_perc

import (
	"github.com/tma15/gostruct"
)

type Node struct {
	Obs   string  /* 観測した文字列 */
	Score float64 /* このノードに付与された重み */
	X     int     /* 観測した系列中での位置  */
	Y     int     /* 隠れ状態のid */
	Fs    []int   /* このノードが持つ素性のリスト */
	LPath []*Path /* ラティス上でこのノードの左側に付いているエッジ */
	//         RPath     []*Path /* ラティス上でこのノードの右側に付いているエッジ */
	Prev      *Node   /* ラティス上で最も良いスコアを持つ左側のノード */
	BestScore float64 /* */
}

func NewNode() Node {
	this := Node{
		LPath: make([]*Path, 0, 10),
		//                 RPath: make([]*Path, 0, 10),
	}
	return this
}

func (this *Node) Features(index gostruct.Index) []string {
	res := make([]string, 0, 10)
	for _, fid := range this.Fs {
		res = append(res, index.Elems[fid])
	}
	return res
}
