package hmm_perc

import (
	"bufio"
	"fmt"
	"github.com/tma15/gostruct"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type FeatureIndex struct {
	NodeFeature gostruct.Index
	NodeWeight  gostruct.Matrix
	//         EdgeFeature gostruct.Index
	EdgeWeight gostruct.Matrix
	Output     gostruct.Index
	Unigrams   []Macro /* Macro defined in template file */
	Bigrams    []Macro /* Macro defined in tempalte file */
	IntRegex   *regexp.Regexp
	Templates  string
}

func NewFeatureIndex() FeatureIndex {
	this := FeatureIndex{
		NodeFeature: gostruct.NewIndex(),
		//                 EdgeFeature: gostruct.NewIndex(),
		Output:   gostruct.NewIndex(),
		IntRegex: regexp.MustCompile(`-?[0-9]+`),
		Unigrams: make([]Macro, 0, 10),
		Bigrams:  make([]Macro, 0, 10),
	}
	return this
}

func FeatureIndexFromFile(model_file, template_file string) FeatureIndex {
	fi := LoadFeatureIndex(model_file)
	fi.IntRegex = regexp.MustCompile(`-?[0-9]+`)
	fi.Unigrams = make([]Macro, 0, 10)
	fi.Bigrams = make([]Macro, 0, 10)
	fi.openTemplate(template_file)
	return fi
}

func (this *FeatureIndex) parseMacros(text string) {
	var m Macro
	var prefix string
	sp := strings.Split(text, ":")
	if len(sp) > 1 {
		prefix = sp[0]
		macros := strings.Split(sp[1], "/")

		cols := make([]int, 0, 10)
		pos := make([]int, 0, 10)
		var n int = 0
		for _, macro := range macros {
			ret := this.IntRegex.FindAllString(macro, -1)
			relpos, _ := strconv.Atoi(ret[0])
			col, _ := strconv.Atoi(ret[1])
			pos = append(pos, relpos)
			cols = append(cols, col)
			n++
		}
		m = NewMacro(prefix, pos, cols, n)
	} else {
		prefix = text
		m = NewMacro(prefix, []int{}, []int{}, 0)
	}
	if string(prefix[0]) == "U" {
		var exists bool = false
		for _, u := range this.Unigrams {
			if u.IsSame(m) {
				exists = true
				continue
			}
		}
		if !exists {
			this.Unigrams = append(this.Unigrams, m)
		}
	} else if string(prefix[0]) == "B" {
		var exists bool = false
		for _, b := range this.Bigrams {
			if b.IsSame(m) {
				exists = true
				continue
			}
		}
		if !exists {
			this.Bigrams = append(this.Bigrams, m)
		}
	}

}

func (this *FeatureIndex) Open(template_file, train_file string) {
	this.openTemplate(template_file)
	this.openTagSet(train_file)
}

func (this *FeatureIndex) openTemplate(template_file string) {
	fp, err := os.Open(template_file)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	var text string
	for scanner.Scan() {
		text = strings.TrimSpace(scanner.Text())
		if len(text) == 0 {
			continue
		}

		/* ignore lines that start with #. These lines are handled as comments  */
		if string(text[0]) == "#" {
			continue
		}
		this.parseMacros(text)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (this *FeatureIndex) Fire(curr int, x [][]string, Ngrams []Macro) []int {
	var id_f int
	fs := make([]int, 0, 10)
	for _, macro := range Ngrams {
		_features := make([]string, 0, macro.N)
		for i := 0; i < macro.N; i++ {
			if curr+macro.Pos[i] >= len(x) {
				_features = append(_features, "__EOS__")
			} else if curr+macro.Pos[i] < 0 {
				_features = append(_features, "__BOS__")
			} else {
				_features = append(_features, x[curr+macro.Pos[i]][macro.Col[i]])
			}
		}
		feature := fmt.Sprintf("%s:%s", macro.Prefix, strings.Join(_features, "/"))
		if string(macro.Prefix[0]) == "U" {
			id_f = this.NodeFeature.GetIdAndAddElemIfNotExists(feature)
			//                 } else if string(macro.Prefix[0]) == "B" {
			//                         id_f = this.EdgeFeature.GetIdAndAddElemIfNotExists(feature)
		}
		fs = append(fs, id_f)
	}
	return fs
}

func (this *FeatureIndex) openTagSet(train_file string) {
	fp, err := os.Open(train_file)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	var text string
	for scanner.Scan() {
		text = strings.TrimSpace(scanner.Text())
		//                 sp := strings.Split(text, "\t")
		sp := strings.Split(text, " ") /* conull format */
		output := sp[len(sp)-1]
		if output == "" {
			continue
		}
		this.Output.GetIdAndAddElemIfNotExists(output)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (this *FeatureIndex) CalcNodeScore(node *Node) {
	y := node.Y
	s := 0.
	for _, fid := range node.Fs {
		this.NodeWeight.Resize(y, fid) /* 素性はラベルと観測の対 */
		s += this.NodeWeight[y][fid]
	}
	node.Score = s
}

func (this *FeatureIndex) CalcPathScore(path *Path) {
	y1 := path.LNode.Y
	y2 := path.RNode.Y
	s := 0.
	//         for _, fid := range path.Fs {
	//                 offset := y1*this.Output.Size() + y2
	//                 this.EdgeWeight.Resize(offset, fid)
	//                 s += this.EdgeWeight[offset][fid]
	//         }
	this.EdgeWeight.Resize(y1, y2)
	s += this.EdgeWeight[y1][y2]
	path.Score = s
}

func (this *FeatureIndex) BuildFeatures(tagger *Tagger) {
	nodes := make([][]Node, 0, 10)
	x := tagger.X

	// 単語ごとに素性idのリストへ変換し、nodeへ付与
	for i := 0; i < len(x); i++ {
		fs := this.Fire(i, x, this.Unigrams)
		nodes = append(nodes, make([]Node, 0, 10))
		for j, _ := range this.Output.Elems {
			node := NewNode()
			node.X = i
			node.Y = j
			node.Fs = fs
			node.Obs = x[i][0]
			this.CalcNodeScore(&node)
			nodes[i] = append(nodes[i], node)
		}
	}

	// node間にエッジ (path) を張る
	for i := 1; i < len(x); i++ {
		for j, _ := range this.Output.Elems {
			for k, _ := range this.Output.Elems {
				path := Path{
					RNode: &nodes[i][k],
					LNode: &nodes[i-1][j],
				}
				this.CalcPathScore(&path)
				path.Add(&nodes[i-1][j], &nodes[i][k])
			}
		}
	}

	tagger.Nodes = nodes
}

func (this *FeatureIndex) Save(model_file string) {
	fp, err := os.OpenFile(model_file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	writer := bufio.NewWriterSize(fp, 4096*32)

	/* outputs */
	writer.WriteString(fmt.Sprintf("%d\n", len(this.Output.Ids)))
	for _, y := range this.Output.Elems {
		writer.WriteString(fmt.Sprintf("%s\n", y))
	}

	/* unigram features */
	writer.WriteString(fmt.Sprintf("%d\n", len(this.NodeFeature.Ids)))
	for _, ft := range this.NodeFeature.Elems {
		writer.WriteString(fmt.Sprintf("%s\n", ft))
	}
	writer.WriteString(fmt.Sprintf("%d\n", len(this.NodeWeight)))
	for y, _ := range this.NodeWeight {
		writer.WriteString(fmt.Sprintf("%d\n", len(this.NodeWeight[y])))
		for _, weight := range this.NodeWeight[y] {
			writer.WriteString(fmt.Sprintf("%f\n", weight))
		}
	}

	/* bigram features */
	//         writer.WriteString(fmt.Sprintf("%d\n", len(this.EdgeFeature.Ids)))
	//         for _, ft := range this.EdgeFeature.Elems {
	//                 writer.WriteString(fmt.Sprintf("%s\n", ft))
	//         }

	/* bigram weight */
	writer.WriteString(fmt.Sprintf("%d\n", len(this.EdgeWeight)))
	for y1, _ := range this.EdgeWeight {
		writer.WriteString(fmt.Sprintf("%d\n", len(this.EdgeWeight[y1])))
		//                 y1_ := this.Output.Elems[y1]
		//                 writer.WriteString(fmt.Sprintf("%d #%s\n", len(this.EdgeWeight[y1]), y1_))
		for _, weight := range this.EdgeWeight[y1] {
			//                         y2_ := this.Output.Elems[y2]
			writer.WriteString(fmt.Sprintf("%f\n", weight))
			//                         writer.WriteString(fmt.Sprintf("%f #%s\n", weight, y2_))
		}
	}
	writer.Flush()
}

func LoadFeatureIndex(model_file string) FeatureIndex {
	var line []byte
	var err error
	fp, err := os.OpenFile(model_file, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	reader := bufio.NewReaderSize(fp, 4096*32)

	/* get output size */
	line, _, err = reader.ReadLine()
	size_output, err := strconv.Atoi(string(line))
	if err != nil {
		panic(err)
	}

	/* read outputs */
	output := gostruct.NewIndex()
	for i := 0; i < size_output; i++ {
		line, _, err = reader.ReadLine()
		if err != nil {
			panic(err)
		}
		output.GetIdAndAddElemIfNotExists(string(line))
	}

	/* read size of unigram features */
	line, _, err = reader.ReadLine()
	size_unigrams, err := strconv.Atoi(string(line))
	if err != nil {
		panic(err)
	}

	/* read unigram features */
	unigrams := gostruct.NewIndex()
	for i := 0; i < size_unigrams; i++ {
		line, _, err = reader.ReadLine()
		if err != nil {
			panic(err)
		}
		unigrams.GetIdAndAddElemIfNotExists(string(line))
	}

	unigram_weight := gostruct.NewMatrix()
	line, _, err = reader.ReadLine()
	if err != nil {
		panic(err)
	}
	size_row, err := strconv.Atoi(string(line))
	/* read weight of unigram features */
	for i := 0; i < size_row; i++ {
		line, _, err = reader.ReadLine()
		if err != nil {
			panic(err)
		}
		size_col, err := strconv.Atoi(string(line))
		for j := 0; j < size_col; j++ {
			line, _, err = reader.ReadLine()
			if err != nil {
				panic(err)
			}
			weight, err := strconv.ParseFloat(string(line), 64)
			if err != nil {
				panic(err)
			}
			unigram_weight.Resize(size_row, size_col)
			unigram_weight[i][j] = weight
		}
	}

	/* read size of bigram features */
	//         line, _, err = reader.ReadLine()
	//         size_bigrams, err := strconv.Atoi(string(line))
	//         if err != nil {
	//                 panic(err)
	//         }

	/* read bigram features */
	//         bigrams := gostruct.NewIndex()
	//         for i := 0; i < size_bigrams; i++ {
	//                 line, _, err = reader.ReadLine()
	//                 if err != nil {
	//                         panic(err)
	//                 }
	//                 bigrams.GetIdAndAddElemIfNotExists(string(line))
	//         }

	bigram_weight := gostruct.NewMatrix()
	line, _, err = reader.ReadLine()
	if err != nil {
		panic(err)
	}
	size_row, err = strconv.Atoi(string(line))
	/* read weight of bigram features */
	for i := 0; i < size_row; i++ {
		line, _, err = reader.ReadLine()
		if err != nil {
			panic(err)
		}
		size_col, err := strconv.Atoi(string(line))
		for j := 0; j < size_col; j++ {
			line, _, err = reader.ReadLine()
			if err != nil {
				panic(err)
			}
			weight, err := strconv.ParseFloat(string(line), 64)
			if err != nil {
				panic(err)
			}
			bigram_weight.Resize(size_row, size_col)
			bigram_weight[i][j] = weight
		}
	}

	fi := FeatureIndex{
		NodeFeature: unigrams,
		//                 EdgeFeature: bigrams,
		NodeWeight: unigram_weight,
		EdgeWeight: bigram_weight,
		Output:     output,
	}
	return fi
}
