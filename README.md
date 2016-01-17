# gostruct
A library of supervised structured output learning such as NP chunking

To use this program, run the folowing command:
```
go get github.com/tma15/gostruct
```

## Supported Algorithms
- Structured Perceptron

## Data Format
Data format of this program is CoNLL 2000 format.

see: [Chunking](http://www.cnts.ua.ac.be/conll2000/chunking/)

## Training
```
cd $GOPATH/src/github.com/tma15/gostruct
wget http://www.cnts.ua.ac.be/conll2000/chunking/train.txt.gz
zcat train.txt.gz > train.txt
./gostruct train -t ./hmm_perc/sample/example.tmp -a hmmperc -m model train.txt
```

## Testing
```
cd $GOPATH/src/github.com/tma15/gostruct
wget http://www.cnts.ua.ac.be/conll2000/chunking/test.txt.gz
zcat train.txt.gz > test.txt
./gostruct test -v -t ./hmm_perc/sample/example.tmp -a hmmperc -m model test.txt
```

## References
- Michael Collins, "Discriminative Training Methods for Hidden Markov Models:
Theory and Experiments with Perceptron Algorithms",  EMNLP, 2002.

## License
see `LICENSE.txt`

## TODO
- Benchmark test
- Refine README
- Implement other online learning algorithms
