---
title:          How I manage all my git repositories with ggman
date:           2025-07-13
author:         Tom Wiesing 
authorLink:     https://tkw01536.de

description:    An overview of why and how I wrote a tool called ggman to manage all my git repositories. 

draft:          true
---

Both at work and in my free time I interact with lots of different git repositories - across my machines I usually have 100 different repositories checked out. 
Maintaining these clones by hand might be possible, but I am way too lazy for that and have written a [go](http://go.dev) program called [ggman](https://github.com/tkw1536/ggman) to manage all of these in a simple way. 

In this post I want to write about how this tool came to be, and give a brief introduction into it's design and usage.   

## Lots of repositories

Sometime around late 2013 when studying at university I was contributing to a research project called [MathHub](https://mathhub.info/) as part of working in a research group. 
As part of the project, the group created various different [git repositories](https://en.wikipedia.org/wiki/Git) on a private [GitLab](https://about.gitlab.com) instance `gl.mathhub.info` that held source code in a [DSL](https://en.wikipedia.org/wiki/Domain-specific_language). 
When working with this content, it was normal for me and others to constantly have somewhere between 20 and 50 repositories organized in different gitlab groups cloned on our machines. 
We then built various other content from these source files. 

While the `git` client could take care of the actual cloning, it soon became clear that manually maintaining these clones by hand was tedious and needed tool support. 
For this purpose we ended up building a tool in Python -- called [LocalMathHub](https://github.com/MathHubInfo/Legacy-localmh) or `lmh` for short -- to maintain a local tree of cloned repositories, `git clone`ing and `git pull`ing them as needed. 
I took over building `lmh` sometime in early 2014. 

`lmh` cloned the various repositories into a structure that mirrored the repository / group structure on the GitLab instance.
For instance the repository `gl.mathhub.info/hello/world` would end up in a folder `$HOME/MathHub/hello/world`, the repository `meta/inf` would end up in a folder `$HOME/MathHub/meta/inf` and so on.
It furthermore supported our domain-specific build processes, to create research data. 
To this end, `lmh` also resolved dependencies between repositories, acting very much like our own package manager similar to [npm](https://www.npmjs.com) or [pip](https://pypi.org/project/pip/).  
Our building processes also required a full [LaTeX](https://en.wikipedia.org/wiki/LaTeX) installation, so the tool ended up being [Docker](https://www.docker.com)ized. 

## Writing GitManager 

Fast forward to a couple years later to 2016 and I was working with lots of different code in lots of different repositories from GitHub. 
I liked a somewhat organized hard disk, so I again wanted to have a local tree mirroring the structure of repositories on GitHub. 
But as GitHub was not our private GitLab instance, `lmh` didn't support it.
Besides, `lmh` was very much focused on the building process and was dockerized - that seemed overkill for the task at hand. 

So I wrote a very simple tool in Python I called [GitManager](https://github.com/tkw1536/GitManager) - that did exactly this. 
It had a configuration file of which repositories to clone where, and supported operations like `git pull` and `git push` across all clones. 
This made my job simpler, but was still annoying to work with as I had to update a configuration file every time I got a new repository. 
Nonetheless I ended up using it for all my clones for a couple years. 

## Design Goals for ggman

Fast forward again a couple more years to 2019 and I got annoyed at needing to update that configuration file. 
So I decided to redo my repository management again. 

I looked around and there were other tools at the time - however these usually had some downsides:
- they were limited to one repository provider (e.g. the [GitHub CLI](http://cli.github.com) only worked with GitHub); or
- They enforced a flat directory structure (e.g everything goes directly under a `Projects` folder); or
- The tool was only available from within an IDE or GUI.

All of these made them unfit for my workflow. 
As I was starting out with [go](https://go.dev) at the time, I decided to use the opportunity to learn the language and start writing my own tool. 
I couldn't think of a good name, and ended up with `ggman` - with "man" standing for "manager" and the gs standing for "git" and "go". 

In order to fit my own workflow, and to prevent me having to rewrite the tool again on the future, I decided that ggman should:

- be command-line first;
- be simple to install, configure and use;
- encourage an obvious hierarchical directory structure, but remain fully functional with any directory structure;
- remain free of provider-specific code; and
- not store any repository-specific data outside of the repositories themselves (enabling the user to switch back to only git at any point).

## the design of ggman

(this paragraph isn't done)

I want to give a brief introduction as to the design of ggman. 
To this end, I feel it is simplest to show how it is used. 

The source code of ggman lives [on GitHub](https://github.com/tkw1536/ggman), the built program consists of a single binary that is dropped into the user's `$PATH` to install.
The binary optionally requires that the user has `git` installed, but will automatically fall back to the [go-git](https://github.com/go-git/go-git) library if not. 

Once installed, ggman manages all git repositories inside a given root directory, and automatically sets up new repositories relative to the URLs they are cloned from. 
This root folder defaults to `~/Projects` but can be customized using the `$GGROOT` environment variable. 

- TODO: `ggman clone`
- TODO: `ggman ls` and filters
- TODO: `ggcd` alias

ggman has lots more functionality, but that feels like a little bit too much for this post. 
I encourage you to have a look at the [README](https://github.com/tkw1536/ggman?tab=readme-ov-file#ggman) or ask me if you're interested. 

## Conclusion

And that is already all I want to say for now, thank you for reading what basically amounts to a wall of text. 

In summary: I work with lots of git repositories.
I wrote a tool called `ggman` to maintain and expand a local directory structure of all of these. 
It can locate where specific repositories are cloned to, and run operations such as `git pull` or `git push` across all of them. 

Feel free to try it out and give me feedback at https://github.com/tkw1536/ggman. 
