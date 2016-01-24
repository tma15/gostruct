package gostruct

import (
	"bufio"
	//         "fmt"
	"io"
	"os"
	"strings"
)

func ReadCoNLLFormat(filename string) ([][][]string, [][]string) {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)

	x := make([][][]string, 0, 1000)
	y := make([][]string, 0, 1000)
	var x_ [][]string = make([][]string, 0, 100)
	var y_ []string = make([]string, 0, 100)
	for scanner.Scan() {
		text := scanner.Text()
		sp := strings.Split(text, " ")
		if len(sp) == 1 {
			x = append(x, x_)
			y = append(y, y_)
			x_ = make([][]string, 0, 100)
			y_ = make([]string, 0, 100)
			continue
		}
		x_i := make([]string, 0, len(sp))
		for i := 0; i < len(sp)-1; i++ {
			x_i = append(x_i, sp[i])
			//                         x_ = append(x_, []string{sp[0], sp[1]})
		}
		x_ = append(x_, x_i)
		//                 x_ = append(x_, []string{sp[0], sp[1]})
		y_ = append(y_, sp[len(sp)-1])
	}
	return x, y
}

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
