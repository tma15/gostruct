package gostruct

type Sequence struct {
	x []string
	y []string
}

func (s *Sequence) Len() int {
	return len(s.x)
}
