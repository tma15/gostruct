package hmm_perc

import (
	"fmt"
	"testing"
)

func TestFitAndPredict(t *testing.T) {
	tagger := NewTagger()
	fi := NewFeatureIndex()
	template_file := "./sample/example.tmp"
	train_file := "./sample/small.txt"
	fi.Open(template_file, train_file)
	tagger.SetFeatureIndex(fi)

	x := [][]string{
		[]string{"Confidence", "B-NP"},
		[]string{"in", "IN"},
		[]string{"the", "DT"},
		[]string{"pound", "NN"},
	}
	y := []string{
		"B-NP",
		"B-PP",
		"B-NP",
		"I-NP",
	}

	for i := 0; i < 2; i++ {
		tagger.Fit(&x, &y)
	}
	pred := tagger.Predict(x)

	fmt.Println("true:", y)
	fmt.Println("pred:", pred)

	for i := 0; i < len(y); i++ {
		if pred[i] != y[i] {
			t.Error("err")
		}
	}
}
