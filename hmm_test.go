package gostruct

import (
	"fmt"
	"testing"
)

func TestSmallHMM(t *testing.T) {
	hmm := NewHMM()

	x := [][]string{
		[]string{
			"a", "b",
		},
		[]string{
			"a", "a",
		},
	}
	y := [][]string{
		[]string{
			"p1", "p2",
		},
		[]string{
			"p1", "p1",
		},
	}

	hmm.Fit(&x, &y)
	num_w := 2
	am := hmm.Forward(num_w, &x[0])
	tags := hmm.Backward(num_w, am)
	fmt.Println(tags)
}

func TestHMM2(t *testing.T) {
	hmm := NewHMM()

	x_train, y_train := LoadTrainFile("/home/makino/code/nlptutorial/data/wiki-en-train.norm_pos")

	hmm.Fit(&x_train, &y_train)
	hmm.Save("model")

	hmm2 := LoadHMM("model")

	x, y := LoadTrainFile("/home/makino/code/nlptutorial/data/wiki-en-test.norm_pos")
	num_data := len(x)
	num_corr := 0.
	num_corr2 := 0.
	num_total := 0.
	num_total2 := 0.
	for i := 0; i < num_data; i++ {
		predy := hmm.Predict(&x[i])
		predy2 := hmm2.Predict(&x[i])
		fmt.Println(y[i])
		fmt.Println(predy2)
		fmt.Println()

		for j := 0; j < len(predy2); j++ {
			if y[i][j] == predy2[j] {
				num_corr++
			}
			num_total++
		}
		for j := 0; j < len(predy2); j++ {
			if predy[j] == predy2[j] {
				num_corr2++
			}
			num_total2++
		}
	}
	fmt.Println(num_corr/num_total, int(num_corr), int(num_total))
	fmt.Println(num_corr2/num_total2, int(num_corr2), int(num_total2))
}
