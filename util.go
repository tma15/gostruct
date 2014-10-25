package main

import (
	"bufio"
//         "fmt"
	"io"
	"os"
	"strings"
)

func LoadTrainFile(fname string) ([][]string, [][]string) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReaderSize(fp, 4096*64)
	X := [][]string{}
	y := [][]string{}

	var pos string
	var word string
	var x_i []string
	var y_i []string
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		sp := strings.Split(string(line), " ")
		x_i = []string{}
		y_i = []string{}
		for _, elem := range sp {
			spelem := strings.Split(elem, "_")
			word = spelem[0]
			pos = spelem[1]
			//                     fmt.Println(word, pos)
			x_i = append(x_i, word)
			y_i = append(y_i, pos)
		}
                X = append(X, x_i)
                y = append(y, y_i)
//                 fmt.Println(x_i)
//                 fmt.Println(y_i)
//                 fmt.Println("")
	}
	return X, y
}

func LoadTestFile(fname string) [][]string {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReaderSize(fp, 4096*64)
	X := [][]string{}

	var word string
	var x_i []string
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		sp := strings.Split(string(line), " ")
		x_i = []string{}
		for i := 0; i < len(sp); i++ {
			word = sp[i]
			x_i = append(x_i, word)
		}
                X = append(X, x_i)
	}
	return X
}

