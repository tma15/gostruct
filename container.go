package gostruct

type Vector []float64

func NewVector() Vector {
	this := Vector{}
	return this
}

func (this *Vector) Resize(size int) {
	if len(*this) <= size {
		for i := len(*this); i <= size; i++ {
			*this = append(*this, 0.)
		}
	}
}

type Matrix [][]float64

func NewMatrix() Matrix {
	this := Matrix{}
	return this
}

func (this *Matrix) Resize(rowId, colId int) {
	if len(*this) <= rowId {
		for i := len(*this); i <= rowId; i++ {
			*this = append(*this, make([]float64, 0, 1000))
		}
	}
	for i, _ := range *this {
		if len((*this)[i]) <= colId {
			for j := len((*this)[i]); j <= colId; j++ {
				(*this)[i] = append((*this)[i], 0.)
			}
		}
	}
}

func (this *Matrix) Fill(val float64) {
	for i := 0; i < len(*this); i++ {
		for j := 0; j < len((*this)[i]); j++ {
			(*this)[i][j] = val
		}
	}
}
