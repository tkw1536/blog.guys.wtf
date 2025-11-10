---
title:          On Typst vs LaTeX 
date:           2025-10-11
author:         Tom Wiesing 
authorLink:     https://tkw01536.de

description:    On Typst and LaTeX

draft:          true
---

I was recently asked by a friend for a quick LaTeX [^1] template on some lecture notes.
I answered that I was using Typst [^2] for any kind of quick notes at this point. 
I then thought I would write a quick post about why I use it. 

## On LaTeX

LaTeX is a typesetting system widely used in academia.
For those not technically versed, in this case "typesetting" refers to producing a paper.
It acts similar to Word and PowerPoint, except that users "program" their documents, feed it through LaTeX, and out comes a `PDF` file.
It's use is widespread: Almost every scientist in mathematics, physics and computer science makes use of it[^4].
If you have ever read a scientific paper, you will most likely recognize the LaTeX fonts intuitively.

LaTeX was first released in 1984.
As such it is showing it's age.
Compiling a multi-page document can take several minutes.
A friend's doctoral thesis takes about 15 - 20 minutes.
On top of that LaTeX frequently produces very cryptic errors such as "Overfull \hbox (Badness 1000000)".

With LaTeX the following is a very typical editing cycle:

1. Writing something
2. Run `pdflatex paper.tex` to start to compile it to pdf
3. Wait for a while
3. Notice that it looks terrible in the pdf
4. Change something and `GOTO 2`

Sometimes the pdf also does not compile at all - however you only notice after a minute of running it.
This cycle of editing is very much not instant.

On top of that, the syntax takes some getting used to.
It is not very readable for either humans or machines.
There are many basic tutorials and plugins for many popular editors.

However, when I said above that people "program" their documents I meant it.
People can - and do - reprogram many basic commands.
For example, the following file is valid (La)TeX [^3]:

```tex
^^5clet~^^5ccatcode~`j0~90 13~`"1~`Y2~77 6jdef ZM1"~`#113jdef
YZXXM1M2"M2iM1YZRR"ppYZVV"QuYZWW"aliYZ::"erYZ55M1M2"aM2M1Y~`@
11Z++"jdefY+jif@"YZ99"j@if"bXg"YY"sXpkYYZ33"luYZ <<M1"jedefjx
"j@if"uR:Y"c5esY"#1YYjxYZ^ ^"iceYZ&&"yeYZ//"SeYZ88"DuYZ;;"s Y
Z--M1M2"M2M1YZ77"e-tneYZ66"inYZzzM1M2"anM2M1YZQQ"O-tcYZOO"NoY
Z44"j@if b&YZ__"j.WYZ22"eYZ00"iYZSS"rY+jj", -St-YZee"!YZ!!"uY
Z=="jparY+jv"s,=Y+j|"-2dXmcYZ``"DY+jw"z2tY"jbf<"-!doj| X2d;-a
nt50 lsY9Y+j.M1 "o X2d -2f-tsM1 -ma5otS m0 Xsm0t=YZAAM1"P:dXc
S2m 6 XSpoM19YZ??"8oYZ**"`2cYZ[["E!YZ]]"CoYZ$$"B!YZBB"8a;co3-
bm5AjvYZCC"-ST2;-SFzocg5ll65BjvYZDD"V5tt-o!S p5ss:-!cl5CjvYZE
E"V6q!' 5S!zl!ojvDYZFF"/x z2sS2;paS7jvEYZGG"/pt2m -yccno;nat-
jvjwFYZHH"Qo -!p2ll5m;!lg7jvGYZII"Ovj,as5tljwjvHYZJJ"*j,o-x2-
-sl!tjwjvIYZKK"Unj| Xbt^n2;6fljwjvJYZLL"?j| -ytmpaXsnt5p;!ls-
jvjwKY+j,M1"2m d-mo6M1;YPXmSj.W A./c-n!dj.o B.T:Xj.tW C.V5tS_
D.V- t6j.W E./xtj.W F./pt-m0j.o G. -5Qj.-vo H.Onj.W I.*-m0j.o
J.-nUj|j.o K.?j|j.o Le+jk")Y("Xtj@if [-^n $SS ]!-hcjjlzs.Yjk4
```

This particular example was constructed by David Carlisle in December 1998.
Despite being somewhat artificial, it does demonstrate that is is not immediately possible for a machine to see the contents of a LaTeX document.

A large real-world dataset of LaTeX documents is held by the Cornell e-Print archive [^4].
Most submissions to the site submit their source code along with pdfs.
I keep seeing machine learning projects that want to make use of all of them.
LaTeX is terrible for this, and so projects either try to use the produced PDFs (which aren't exactly nice to work with) or fall back to using paper metadata only [^5].

## On Typst

(to be written)

[^1]: https://en.wikipedia.org/wiki/LaTeX
[^2]: https://typst.app
[^3]: https://ctan.org/tex-archive/macros/plain/contrib/xii-lat
[^4]: https://arxiv.org
[^5]: Side note: What doesn't seem to be widely known is that arXiv has offered https://ar5iv.labs.arxiv.org for a while now.
  It is powered by a software called LaTeXML (which I've worked on in the past) to display submissions as responsive HTML5.
  HTML5 is much nicer to work with than either the (La)TeX sources or pdfs. 
  They also offer a dataset of all papers of all of arxiv in html for research purposes - https://sigmathling.kwarc.info/resources/ar5iv-dataset-2024/. 
  But somehow I've barely seen any projects using this.