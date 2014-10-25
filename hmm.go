package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	BOS = "<s>"
	EOS = "</s>"
)

type HMM struct {
	e      map[string]map[string]float64
	t      map[string]map[string]float64
	p      []string
	lambda float64
}

func NewHMM(l float64) HMM {
	hmm := HMM{map[string]map[string]float64{}, map[string]map[string]float64{}, []string{}, l}
	return hmm
}

func (h *HMM) Fit(X, y [][]string) {
	var xn, ym, yn string
	e := map[string]map[string]float64{}
	t := map[string]map[string]float64{}
	sumy := map[string]float64{}
	sumx := map[string]float64{}

	for i, _ := range X {
		seq := Sequence{X[i], y[i]}
		for j := 0; j <= seq.Len(); j++ {
			if j == 0 {
				ym = BOS
				yn = seq.y[j]
			} else if j == seq.Len() {
				ym = seq.y[j-1]
				yn = EOS
			} else {
				xn = seq.x[j]
				ym = seq.y[j-1]
				yn = seq.y[j]
				if _, ok := e[yn]; ok {
					if _, ok1 := e[yn][xn]; ok1 {
						e[yn][xn] += 1.
					} else {
						e[yn][xn] = 1.
					}
				} else {
					e[yn] = map[string]float64{}
					e[yn][xn] = 1.
				}

				if _, ok := sumx[xn]; ok {
					sumx[xn] += 1.
				} else {
					sumx[xn] = 1.
				}

			}
			if _, ok := t[ym]; ok {
				if _, ok1 := t[ym][yn]; ok1 {
					t[ym][yn] += 1.
				} else {
					t[ym][yn] = 1.
				}
			} else {
				t[ym] = map[string]float64{}
				t[ym][yn] = 1.
			}

			if _, ok := sumy[yn]; ok {
				sumy[yn] += 1.
			} else {
				sumy[yn] = 1.
			}

		}
	}
	for pos, _ := range sumy {
		h.p = append(h.p, pos)
	}
	h.p = append(h.p, BOS)

	for m, _ := range e {
		for n, cnt := range e[m] {
			if _, ok := h.e[m]; !ok {
				h.e[m] = map[string]float64{}
			}
			h.e[m][n] = cnt / sumx[n]
			//             fmt.Println("Pr_e", "(", m, "->", n, ")", cnt, "/", sumx[n], "=", h.e[m][n])
		}
	}

	for m, _ := range t {
		for n, cnt := range t[m] {
			if _, ok := h.t[m]; !ok {
				h.t[m] = map[string]float64{}
			}
			h.t[m][n] = cnt / sumy[n]
			//             fmt.Println("Pr_t", "(", m, "->", n, ")", cnt, "/", sumy[n])
		}
	}
}

func (h *HMM) Forward(num_w int, x []string) map[int]map[string]string {
	bestScore := map[int]map[string]float64{}
	bestScore[0] = map[string]float64{
		BOS: 0.,
	}
	bestEdge := map[int]map[string]string{}
	bestEdge[0] = map[string]string{
		BOS: "",
	}

	var pe, score float64
	var ok, ok1, ok2 bool
	for i := 0; i < num_w; i++ {
		if _, ok := bestScore[i+1]; !ok {
			bestScore[i+1] = map[string]float64{}
		}
		if _, ok := bestEdge[i+1]; !ok {
			bestEdge[i+1] = map[string]string{}
		}
		for _, prev := range h.p {
			for _, next := range h.p {
				_, ok1 = bestScore[i][prev]
				_, ok2 = h.t[prev][next]
				if ok1 && ok2 {
					pe = h.lambda*h.e[next][x[i]] + (1 - h.lambda) + 1./1000. // Smoothing
					score = bestScore[i][prev] - math.Log(h.t[prev][next]) - math.Log(pe)
					_, ok = bestScore[i+1][next]
					if !ok || bestScore[i+1][next] > score {
						bestScore[i+1][next] = score
						bestEdge[i+1][next] = prev
					}
				}
			}
		}
	}

	i := num_w
	if _, ok := bestScore[i+1]; !ok {
		bestScore[i+1] = map[string]float64{}
	}
	if _, ok := bestEdge[i+1]; !ok {
		bestEdge[i+1] = map[string]string{}
	}
	for _, prev := range h.p {
		_, ok1 = bestScore[i][prev]
		_, ok2 = h.t[prev][EOS]
		if ok1 && ok2 {
			score = bestScore[i][prev] - math.Log(h.t[prev][EOS])
			_, ok = bestScore[i+1][EOS]
			if !ok || bestScore[i+1][EOS] > score {
				bestScore[i+1][EOS] = score
				bestEdge[i+1][EOS] = prev
			}
		}
	}
	return bestEdge

}

func (h *HMM) Backward(num_w int, bestEdge map[int]map[string]string) []string {
	tags := []string{}
	nextEdge := bestEdge[num_w+1][EOS]
	for i := num_w; i > 0; i-- {
		tags = append(tags, nextEdge)
		nextEdge = bestEdge[i][nextEdge]
	}
	reversed := []string{}
	for i := len(tags) - 1; i >= 0; i-- {
		reversed = append(reversed, tags[i])
	}
	return reversed
}

func (h *HMM) Predict(x []string) []string {
	num_w := len(x)
	bestEdge := h.Forward(num_w, x)
	return h.Backward(num_w, bestEdge)
}

func SaveHMM(h HMM, fname string) {
	os.Remove(fname)
	modelf, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0777)
	defer modelf.Close()
	if err != nil {
		panic("Faild to open model file")
	}
	writer := bufio.NewWriterSize(modelf, 4096*32)
	writer.WriteString("LAMBDA\n")
	writer.WriteString(fmt.Sprintf("%f\n", h.lambda))
	writer.WriteString("LABEL\n")
	for i := 0; i < len(h.p); i++ {
		writer.WriteString(fmt.Sprintf("%s\n", h.p[i]))
	}

	writer.WriteString("EMISSION\n")
	for y, _ := range h.e {
		for x, pr_e := range h.e[y] {
			writer.WriteString(fmt.Sprintf("%s\t%s\t%f\n", y, x, pr_e))
		}
	}
	writer.WriteString("TRANSITION\n")
	for prev, _ := range h.t {
		for next, pr_t := range h.t[prev] {
			writer.WriteString(fmt.Sprintf("%s\t%s\t%f\n", prev, next, pr_t))
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
	var text, m, prev, next, y, x string
	var lambda, p float64
	labels := []string{}
	e := map[string]map[string]float64{}
	t := map[string]map[string]float64{}
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
		if text == "EMISSION" {
			m = text
			continue
		}
		if text == "TRANSITION" {
			m = text
			continue
		}
		if m == "LAMBDA" {
			lambda, _ = strconv.ParseFloat(text, 64)
		}
		if m == "LABEL" {
			labels = append(labels, text)
		}
		if m == "TRANSITION" {
			sp := strings.Split(text, "\t")
			prev = sp[0]
			next = sp[1]
			p, _ = strconv.ParseFloat(sp[2], 64)
			if _, ok := t[prev]; !ok {
				t[prev] = map[string]float64{}
			}
			t[prev][next] = p
		}
		if m == "EMISSION" {
			sp := strings.Split(text, "\t")
			y = sp[0]
			x = sp[1]
			p, _ = strconv.ParseFloat(sp[2], 64)
			if _, ok := e[y]; !ok {
				e[y] = map[string]float64{}
			}
			e[y][x] = p
		}
	}
	h := HMM{e, t, labels, lambda}
	return h
}

type Sequence struct {
	x []string
	y []string
}

func (s *Sequence) Len() int {
	return len(s.x)
}

func main() {
    flag.Usage = func () {
        flag.PrintDefaults()
    }
    fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    var (
        modelfile = fs.String("m", "model", "model file")
        input = fs.String("i", "input", "input")
        lambda = fs.Float64("l", 0.9, "smoothing parameter")
    )
    if os.Args[1] == "-h" {
        fmt.Println("./hmm [train|test] OPTIONS")
       os.Exit(1)
    }

    mode := os.Args[1]
    fs.Parse(os.Args[2:])
    fmt.Println("mode:", mode)
    fmt.Println("modelfile:", *modelfile)
    fmt.Println("input:", *input)

    if mode == "train" {
	X, y := LoadTrainFile(*input)
	h := NewHMM(*lambda)
	h.Fit(X, y)
	SaveHMM(h, *modelfile)
    } else if mode == "test" {
	X := LoadTestFile(*input)
	h := LoadHMM(*modelfile)
	var y_pred []string
        for i := 0; i < len(X); i++{
//             fmt.Println(X[i])
            y_pred = h.Predict(X[i])
            fmt.Println(strings.Join(y_pred, " "))
        }
    }
}
