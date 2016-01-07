package main

import (
	"flag"
	"fmt"
	"github.com/tma15/gostruct"
	"github.com/tma15/gostruct/hmm_perc"
	"os"
	"strings"
)

func train(args []string) {
	fs := flag.NewFlagSet("train", flag.ExitOnError)
	var (
		model_file    = fs.String("m", "model", "model file")
		input         = fs.String("i", "input", "input")
		algorithm     = fs.String("a", "algorithm", "algorithm")
		template_file = fs.String("t", "../hmm_perc/sample/example.tmp", "template file")
	)

	fs.Parse(args)

	switch *algorithm {
	case "hmm":
		x, y := gostruct.LoadTrainFile(*input)
		h := gostruct.NewHMM()
		h.Fit(&x, &y)
		h.Save(*model_file)
	case "hmmperc":
		feature_index := hmm_perc.NewFeatureIndex()
		for _, train_file := range fs.Args() {
			feature_index.Open(*template_file, train_file)
		}

		// train
		tagger := hmm_perc.NewTagger()
		tagger.SetFeatureIndex(feature_index)

		for _, train_file := range fs.Args() {
			x, y := gostruct.ReadCoNLLFormat(train_file)
			n := len(x)
			fmt.Println(n)
			for i := 0; i < n; i++ {
				tagger.Fit(&x[i], &y[i])
			}
		}
		tagger.Save(*model_file)

	}

}

func test(args []string) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	var (
		model_file    = fs.String("m", "model", "model file")
		template_file = fs.String("t", "../hmm_perc/sample/example.tmp", "template file")
		input         = fs.String("i", "input", "input")
		algorithm     = fs.String("a", "algorithm", "algorithm")
	)
	fs.Parse(args)

	switch *algorithm {
	case "hmm":
		X := gostruct.LoadTestFile(*input)
		h := gostruct.LoadHMM(*model_file)
		var y_pred []string
		for i := 0; i < len(X); i++ {
			y_pred = h.Predict(&X[i])
			fmt.Println(strings.Join(y_pred, " "))
		}
	case "hmmperc":
		fi := hmm_perc.FeatureIndexFromFile(*model_file, *template_file)
		fmt.Println(len(fi.Output.Ids))

		tagger := hmm_perc.NewTagger()
		tagger.SetFeatureIndex(fi)

		// test
		var token_match int = 0
		var token_total int = 0
		var sent_match int = 0
		var sent_total int = 0

		for _, test_file := range fs.Args() {
			x, y := gostruct.ReadCoNLLFormat(test_file)
			n := len(x)
			for i := 0; i < n; i++ {
				pred := tagger.Predict(x[i])
				//                                 fmt.Println(x[i])
				//                                 fmt.Println("true", y[i])
				//                                 fmt.Println("pred", pred)
				//                                 fmt.Println()

				for j := 0; j < len(y[i]); j++ {
					if y[i][j] == pred[j] {
						token_match++
					}
					token_total++
				}

				if hmm_perc.IsTheSame(y[i], pred) {
					sent_match++
				}

				sent_total++

				fmt.Println(fmt.Sprintf("token:%f (%d/%d)",
					float64(token_match)/float64(token_total), token_match, token_total))

				//                                 fmt.Println(fmt.Sprintf("sent:%f (%d/%d)",
				//                                         float64(sent_match)/float64(sent_total), sent_match, sent_total))
			}
		}
		fmt.Println(fmt.Sprintf("token:%f (%d/%d)",
			float64(token_match)/float64(token_total), token_match, token_total))

	}

}

func main() {
	var usage = `
Usage of %s <Command> [Options]

Commands:
  train
  test

`

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "train":
		train(args[1:])
	case "test":
		test(args[1:])
	default:
		flag.Usage()
		os.Exit(1)
	}
}
