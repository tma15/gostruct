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
			for i := 0; i < n; i++ {
				tagger.Fit(&x[i], &y[i])
			}
		}
		tagger.Save(*model_file)

		// closed test
		//                 var token_match int = 0
		//                 var token_total int = 0
		//                 var sent_match int = 0
		//                 var sent_total int = 0

		//                 fi := hmm_perc.FeatureIndexFromFile(*model_file, *template_file)

		//                 tagger2 := hmm_perc.NewTagger()
		//                 tagger2.SetFeatureIndex(fi)

		//                 for _, train_file := range fs.Args() {
		//                         x, y := gostruct.ReadCoNLLFormat(train_file)
		//                         n := len(x)
		//                         for i := 0; i < n; i++ {
		//                                 pred := tagger.Predict(x[i])
		//                                 pred2 := tagger.Predict(x[i])

		//                                 for j := 0; j < len(y[i]); j++ {
		//                                         y_true := y[i][j]
		//                                         y_pred := pred[j]
		//                                         if y_true == y_pred {
		//                                                 token_match++
		//                                         }
		//                                         token_total++
		//                                 }

		//                                 if hmm_perc.IsTheSame(y[i], pred) {
		//                                         sent_match++
		//                                 }
		//                                 if !hmm_perc.IsTheSame(pred, pred2) {
		//                                         fmt.Println("pred1", pred)
		//                                         fmt.Println("pred2", pred2)
		//                                         fmt.Println("==")
		//                                 }

		//                                 sent_total++
		//                         }
		//                 }
		//                 fmt.Println(fmt.Sprintf("accuracy:%f (%d/%d)",
		//                         float64(token_match)/float64(token_total), token_match, token_total))

	}

}

func test(args []string) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	var (
		model_file    = fs.String("m", "model", "model file")
		template_file = fs.String("t", "../hmm_perc/sample/example.tmp", "template file")
		input         = fs.String("i", "", "input")
		algorithm     = fs.String("a", "", "algorithm")
		verbose       = fs.Bool("v", false, "verbose mode")
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

		tagger := hmm_perc.NewTagger()
		tagger.SetFeatureIndex(fi)

		// test
		var token_match int = 0
		var token_total int = 0
		var sent_match int = 0
		var sent_total int = 0

		// 正解のラベルごとの予測ラベルの分布
		// precision: 予測ラベル == 真のラベル / 予測ラベル
		// recall: 予測ラベル == 真のラベル / 真のラベル
		true_dist := make(map[string]float64)
		true_sum := make(map[string]float64)
		pred_sum := make(map[string]float64)

		for _, test_file := range fs.Args() {
			x, y := gostruct.ReadCoNLLFormat(test_file)
			n := len(x)
			for i := 0; i < n; i++ {
				pred := tagger.Predict(x[i])

				for j := 0; j < len(y[i]); j++ {
					y_true := y[i][j]
					y_pred := pred[j]
					if y_true == y_pred {
						token_match++
						true_dist[y_true] += 1.
					}
					token_total++
					true_sum[y_true] += 1.
					pred_sum[y_pred] += 1.
				}

				if hmm_perc.IsTheSame(y[i], pred) {
					sent_match++
				}

				sent_total++

				if *verbose {
					for t := 0; t < len(x[i]); t++ {
						fmt.Println(fmt.Sprintf("%s\t%s\t%s", x[i][t][0], y[i][t], pred[t]))
					}
				}
			}
		}
		//                 fmt.Println(fmt.Sprintf("accuracy:%f (%d/%d)",
		//                         float64(token_match)/float64(token_total), token_match, token_total))

		tp := 0.
		num_pred := 0.
		num_true := 0.
		for y, _ := range true_dist {
			tp += true_dist[y]

			prec := true_dist[y] / pred_sum[y]
			num_pred += pred_sum[y]
			prectext := fmt.Sprintf("%.3f (%d/%d)", prec,
				int(true_dist[y]), int(pred_sum[y]))

			reca := true_dist[y] / true_sum[y]
			num_true += true_sum[y]
			recatext := fmt.Sprintf("%.3f (%d/%d)", reca,
				int(true_dist[y]), int(true_sum[y]))

			fmt.Println(fmt.Sprintf("%10s %10s %10s", y, prectext, recatext))
		}
		prec_avg := tp / num_pred
		reca_avg := tp / num_true
		f_score := (2. * prec_avg * reca_avg) / (prec_avg + reca_avg)

		fmt.Println("")
		fmt.Println(fmt.Sprintf("precision: %f (%d/%d)", prec_avg,
			int(tp), int(num_pred)))
		fmt.Println(fmt.Sprintf("recall: %f (%d/%d)", reca_avg,
			int(tp), int(num_true)))
		fmt.Println(fmt.Sprintf("f-score: %f", f_score))
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
