package main

import (
	"flag"
	"fmt"
	"github.com/tma15/gostruct"
	"os"
	"strings"
)

func train(args []string) {
	fs := flag.NewFlagSet("train", flag.ExitOnError)
	var (
		modelfile = fs.String("m", "model", "model file")
		input     = fs.String("i", "input", "input")
	)

	fs.Parse(args)

	x, y := gostruct.LoadTrainFile(*input)
	h := gostruct.NewHMM()
	h.Fit(&x, &y)
	h.Save(*modelfile)
}

func test(args []string) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	var (
		modelfile = fs.String("m", "model", "model file")
		input     = fs.String("i", "input", "input")
	)
	fs.Parse(args)

	X := gostruct.LoadTestFile(*input)
	h := gostruct.LoadHMM(*modelfile)
	var y_pred []string
	for i := 0; i < len(X); i++ {
		y_pred = h.Predict(&X[i])
		fmt.Println(strings.Join(y_pred, " "))
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
