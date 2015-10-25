package gostruct

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type HMM struct {
	XIndex        Index
	YIndex        Index
	SumEmission   Vector /* number of emission of x */
	SumTransition Vector /* number of transition of y */
	TrToEOS       Vector
	TrFromBOS     Vector
	Emission      Matrix
	Transition    Matrix
}

func NewHMM(l float64) HMM {
	hmm := HMM{
		XIndex:        NewIndex(),
		YIndex:        NewIndex(),
		SumEmission:   NewVector(),
		SumTransition: NewVector(),
		Emission:      NewMatrix(),
		Transition:    NewMatrix(),
		TrToEOS:       NewVector(),
		TrFromBOS:     NewVector(),
	}
	return hmm
}

func (h *HMM) Fit(X, y [][]string) {
	var xn, ym, yn string
	var iym, iyn, ixn int

	num_data := float64(len(X))
	for i, _ := range X {
		seq := Sequence{X[i], y[i]}
		for j := 0; j <= seq.Len(); j++ {
			if j == 0 {
				/* BOSと先頭の品詞 */
				ym = BOS
				yn = seq.y[j]
				iyn = h.YIndex.GetIdAndAddElemIfNotExists(yn)
				xn = seq.x[j]
				ixn = h.XIndex.GetIdAndAddElemIfNotExists(xn)
				h.Emission.Resize(iyn, ixn)
				h.Emission[iyn][ixn] += 1.
				h.TrFromBOS.Resize(iyn)
				h.TrFromBOS[iyn] += 1. / num_data
			} else if j == seq.Len() {
				/* 末尾の品詞とEOS */
				ym = seq.y[j-1]
				yn = EOS
				iym = h.YIndex.GetIdAndAddElemIfNotExists(ym)
				h.TrToEOS.Resize(iym)
				h.TrToEOS[iym] += 1. / num_data
			} else {
				xn = seq.x[j]
				ym = seq.y[j-1]
				yn = seq.y[j]
				iym = h.YIndex.GetIdAndAddElemIfNotExists(ym)
				iyn = h.YIndex.GetIdAndAddElemIfNotExists(yn)
				ixn = h.XIndex.GetIdAndAddElemIfNotExists(xn)

				h.Emission.Resize(iyn, ixn)
				h.Emission[iyn][ixn] += 1.

				h.SumEmission.Resize(ixn)
				h.SumEmission[ixn] += 1.

				/* iyn > iym。(iyn, iyn)の正方行列にresize */
				h.Transition.Resize(iyn, iyn)
				h.Transition[iym][iyn] += 1.
				h.SumTransition.Resize(iyn)
				h.SumTransition[iyn] += 1.
			}
		}
	}
	for y := range h.Emission {
		for x := range h.Emission[y] {
			if h.SumEmission[x] > 0 {
				h.Emission[y][x] /= float64(h.SumEmission[x])
			}
		}
	}
	for ym := range h.Transition {
		for yn := range h.Transition[ym] {
			if h.SumTransition[yn] > 0 {
				h.Transition[ym][yn] /= float64(h.SumTransition[yn])
			}
		}
	}
}

func (h *HMM) nl(val float64) float64 {
	if val == 0. {
		return -math.Log(1e-10)
	}
	return -math.Log(val)
}

func (h *HMM) Forward(num_w int, x []string) Matrix {
	ysize := len(h.YIndex.Ids)
	be := NewMatrix()
	am := NewMatrix()
	be.Resize(num_w, ysize-1)
	be.Fill(1e+10)
	am.Resize(num_w, ysize-1)

	/* 初期化 */
	var w int = -1
	var has_w bool
	has_w = h.XIndex.HasElem(x[0])
	if has_w {
		w = h.XIndex.GetId(x[0])
	}
	for i, _ := range h.YIndex.Elems {
		be[0][i] = 0.
		/* BOSからの遷移確率 + iの先頭の単語の放出確率  */
		if i < len(h.TrFromBOS) {
			be[0][i] += h.nl(h.TrFromBOS[i])
		}
		if has_w {
			be[0][i] += h.nl(h.Emission[i][w])
		}
	}

	var hypo float64
	for i := 0; i < num_w-1; i++ {
		has_w = h.XIndex.HasElem(x[i+1])
		if has_w {
			w = h.XIndex.GetId(x[i+1])
		}
		for c, _ := range h.YIndex.Elems {
			for n, _ := range h.YIndex.Elems {
				hypo = be[i][c] + h.nl(h.Transition[c][n])
				if has_w && w < len(h.Emission[n]) {
					hypo += h.nl(h.Emission[n][w])
				}
				if be[i+1][n] > hypo {
					be[i+1][n] = hypo
					/* nextに繋がる最大スコアのエッジはcから出る */
					am[i+1][n] = float64(c)
				}
			}
		}
	}
	for c, _ := range h.YIndex.Elems {
		hypo = 0.
		if c < len(h.TrToEOS) {
			hypo += h.nl(h.TrToEOS[c])
		}
		hypo += be[num_w-1][c]
		if be[num_w][0] > hypo {
			be[num_w][0] = hypo
			am[num_w][0] = float64(c)
		}
	}
	//         for _, row := range be {
	//                 fmt.Println(row)
	//         }
	//         for _, row := range am {
	//                 fmt.Println(row)
	//         }

	return am
}

func (h *HMM) Backward(num_w int, am Matrix) []string {
	tags := make([]string, num_w, num_w)
	argmax := int(am[num_w][0])
	tags[len(tags)-1] = h.YIndex.Elems[argmax]
	for i := 1; i <= num_w-1; i++ {
		argmax = int(am[num_w-i][argmax])
		tags[len(tags)-1-i] = h.YIndex.Elems[argmax]
	}
	return tags
}

func (h *HMM) Predict(x []string) []string {
	num_w := len(x)
	bestEdge := h.Forward(num_w, x)
	return h.Backward(num_w, bestEdge)
}

func (h *HMM) Save(fname string) {
	os.Remove(fname)
	modelf, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0777)
	defer modelf.Close()
	if err != nil {
		panic("Faild to open model file")
	}
	writer := bufio.NewWriterSize(modelf, 4096*32)
	writer.WriteString("LABEL\n")
	for _, elem := range h.YIndex.Elems {
		writer.WriteString(fmt.Sprintf("%s\n", elem))
	}
	writer.WriteString("OBSERVATION\n")
	for _, elem := range h.XIndex.Elems {
		writer.WriteString(fmt.Sprintf("%s\n", elem))
	}

	writer.WriteString("EMISSION\n")
	for y := range h.Emission {
		for x, pr_e := range h.Emission[y] {
			if pr_e > 0 {
				writer.WriteString(fmt.Sprintf("%d\t%d\t%f\n", y, x, pr_e))
			}
		}
	}
	writer.WriteString("TRANSITION\n")
	for prev := range h.Transition {
		for next, pr_t := range h.Transition[prev] {
			if pr_t > 0 {
				writer.WriteString(fmt.Sprintf("%d\t%d\t%f\n", prev, next, pr_t))
			}
		}
	}
	writer.WriteString("TRANSITION_FROM_BOS\n")
	for y, pr_t := range h.TrFromBOS {
		if pr_t > 0 {
			writer.WriteString(fmt.Sprintf("%d\t%f\n", y, pr_t))
		}
	}
	writer.WriteString("TRANSITION_TO_EOS\n")
	for y, pr_t := range h.TrToEOS {
		if pr_t > 0 {
			writer.WriteString(fmt.Sprintf("%d\t%f\n", y, pr_t))
		}
	}
	writer.Flush()
}

func LoadHMM(fname string) HMM {
	modelf, err := os.OpenFile(fname, os.O_RDONLY, 0644)
	if err != nil {
		panic("Failed to load model file")
	}
	reader := bufio.NewReaderSize(modelf, 4096*32)
	var xindex, yindex Index = NewIndex(), NewIndex()
	var text, m string
	var prev, next, y, x int
	var e, t Matrix = NewMatrix(), NewMatrix()
	var bos, eos Vector = NewVector(), NewVector()
	var p float64
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		text = string(line)
		if text == "LAMBDA" {
			m = text
			continue
		}
		if text == "LABEL" {
			m = text
			continue
		}
		if text == "OBSERVATION" {
			m = text
			continue
		}
		if text == "EMISSION" {
			m = text
			continue
		}
		if text == "TRANSITION" {
			m = text
			continue
		}
		if text == "TRANSITION_FROM_BOS" {
			m = text
			continue
		}
		if text == "TRANSITION_TO_EOS" {
			m = text
			continue
		}
		if m == "LABEL" {
			yindex.AddElem(text)
		}
		if m == "OBSERVATION" {
			xindex.AddElem(text)
		}
		if m == "TRANSITION" {
			sp := strings.Split(text, "\t")
			prev, _ = strconv.Atoi(sp[0])
			next, _ = strconv.Atoi(sp[1])
			p, _ = strconv.ParseFloat(sp[2], 64)
			if prev > next {
				t.Resize(prev, prev)
			} else {
				t.Resize(next, next)
			}
			t[prev][next] = p
		}
		if m == "TRANSITION_FROM_BOS" {
			sp := strings.Split(text, "\t")
			y, _ = strconv.Atoi(sp[0])
			p, _ = strconv.ParseFloat(sp[1], 64)
			bos.Resize(y)
			bos[y] = p
		}
		if m == "TRANSITION_TO_BOS" {
			sp := strings.Split(text, "\t")
			y, _ = strconv.Atoi(sp[0])
			p, _ = strconv.ParseFloat(sp[1], 64)
			eos.Resize(y)
			eos[y] = p
		}
		if m == "EMISSION" {
			sp := strings.Split(text, "\t")
			y, _ = strconv.Atoi(sp[0])
			x, _ = strconv.Atoi(sp[1])
			p, _ = strconv.ParseFloat(sp[2], 64)
			e.Resize(y, x)
			e[y][x] = p
		}
	}
	h := HMM{
		Emission:   e,
		Transition: t,
		XIndex:     xindex,
		YIndex:     yindex,
		TrFromBOS:  bos,
		TrToEOS:    eos,
	}
	return h
}
