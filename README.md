gostruct
==
A library of supervised structured learning

Usage
==
```
func main() {
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

  test := []string{
      "a", "b",
  }
  pred := hmm.Predict(&test)
}
```

Example
==
This example uses POS tagging in  [nlptutorial](https://github.com/neubig/nlptutorial).

```
$./gostruct train -i ~/code/nlptutorial/data/wiki-en-train.norm_pos
$./gostruct test -i ~/code/nlptutorial/data/wiki-en-test.norm -m model > my_answer.pos
$~/code/nlptutorial/script/gradepos.pl ~/code/nlptutorial/data/wiki-en-test.pos my_answer.pos3
```

--
Format
--

Training data
---
```
<word>_<POS>  <word>_<POS> ...
```

Test data
---
```
<word> <word> ...
```
