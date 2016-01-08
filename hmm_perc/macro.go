package hmm_perc

func is_same_slice(x, y []int) bool {
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

type Macro struct {
	Prefix string
	Pos    []int
	Col    []int
	N      int // number of words
}

func NewMacro(prefix string, pos, col []int, n int) Macro {
	this := Macro{
		Prefix: prefix,
		Pos:    pos,
		Col:    col,
		N:      n,
	}
	return this
}

func (this *Macro) IsSame(other Macro) bool {
	if this.Prefix == other.Prefix &&
		is_same_slice(this.Pos, other.Pos) &&
		is_same_slice(this.Col, other.Col) &&
		this.N == other.N {
		return true
	}
	return false
}
