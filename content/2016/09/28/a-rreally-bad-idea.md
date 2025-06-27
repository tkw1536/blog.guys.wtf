---
title:      A Rreally bad idea
date:       2016-09-28
author:     Tom Wiesing 
authorLink: https://tkw01536.de
---

[R](https://www.r-project.org/) is a tool for statistical modelers and social scientists to be able to make statistical computations. As such I have come across it in the context of statistical modelling. I consider myself an experienced programmer and when I was first thrown at R I had quite some problems with it.

It is significantly different from any other language I have worked with. In fact it is the only language I know of if I need a variable x with the value ```42``` and I write ```x = 42``` people will look at me funny. In practice, it is only done with ```x <- 42```. Both statements do the exact same thing, but only the second one is used.

You may ask why. R started out with just having the ```<-``` operator, but people asked to use equality instead. Problem was this was already used for equality checking. Eventually the change was implemented and the language now needs to infer from context which equality the user refers to – assignment or checking. This works well in strictly typed functional languages (such as SML or OCaml) but is difficult with R. It seems that it is discouraged to not confuse users, but then why implement it in the first place?

## Packages pollute the global namespace

However this is not the only problem. The biggest fault I can see in R is namespacing — there is not really any. R has a package system that allows installing and loading of packages that extend it with further functionality. A package can be loaded by just

```R
# Load the car package
package(car)
```

Unfortunately this imports any function into the global namespace. Fortunate enough this particular package creates only functions that start with ```cr.```, however this is just a convention and just serves as an identifier. There is no language construct enforcing this and an object ```cr``` simply does not exist.

In principle, if you load a new package, it could do anything to your namespace. That includes overwriting existing values, adding new (potentially big) datatsets and adding new functions with weird names. Before loading a package, there is no real way of knowing what exactly will be loaded and what will be there afterwards.

You might be saying this saying this is not too bad. Hold on, because this is not the worst. Loading an existing dataset into R becomes even worse.

## Datasets out of nowhere

In one of the recent investigations I was doing, I was asked to load an existing dataset into R. This was stored in a so-called RData file. I was also told I could use load(filename) to achieve this. At first I tried the following:

```R
homes <- load("./OregonHomes.Rdata")
```

To my surprise, the variable homes had the string value ```"homes"``` afterwards. This confused the hell out of me – that was certainly not the dataset I was looking to investigate. After some time, I tried to look at the output of just the load call. Oddly enough, this produced no output at all. There was no result, no error message, no nothing.

Fortunately I was using a graphical IDE which also showed me all the variables in the local namespace. I noticed that the variable homes had suddenly changed its value – the dataset that I was looking for.

## So what the hell happened?

Turns out that an Rdata file does not contain a single variable, it contains an entire environment. It is actually just a dictionary from variable names to values. By default the ```load()``` call behaves just like Pythons ```from package import *```, dumping everything in the global namespace.

This is something which should never be done because of a lot of unexpected side effects – accidentally overwriting a variable for example. Importing datasets is thus just as bad as package loading. There seems to be no easy way of just exporting a single variable and loading it later.

Looking at the documentation of the load function now, there is a way to namespace this, but it is non-obvious why one can not just assign it to a namespace by storing the output of the load call. It is also very difficult to just load a single variable.

It seems that this is the behavior of most of Rs code base – it seems more of a philosophy. Everything assumes the global namespace by default, everything else needs to be given by separate arguments. Of course none of this is taught in tutorials or classes – they just leave you to deal with the shitty defaults.

In conclusion, no the title of this post is not a typo – R is just a rreally bad idea when trying to structurally write code.

![Or you can use R, but it's not a real programming language](/media/rmeme.jpg)
