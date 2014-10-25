Usage
==
```
$./hmm train -i ../nlp-programming/data/wiki-en-train.norm_pos
$./hmm test -i ../nlp-programming/data/wiki-en-test.norm > my_answer.pos
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

