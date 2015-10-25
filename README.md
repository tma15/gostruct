gostruct
==
a library of structured learning

Example
==
This example uses [nlptutorial](https://github.com/neubig/nlptutorial).

```
$./gostruct train -i ~/code/nlptutorial/data/wiki-en-train.norm_pos
$./gostruct test -i ~/code/nlptutorial/data/wiki-en-test.norm -m model > my_answer.pos
$~/code/nlptutorial/script/gradepos.pl ~/code/nlptutorial/data/wiki-en-test.pos my_answer.pos3
```

Format
==
---
Training data
---
```
<word>_<POS>  <word>_<POS> ...
```

---
Test data
---
```
<word> <word> ...
```

