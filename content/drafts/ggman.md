---
title:          How I manage all my git repositories with ggman
date:           2025-07-13
author:         Tom Wiesing 
authorLink:     https://tkw01536.de

description:    An overview of why and how I wrote a tool called ggman to manage all my git repositories. 

draft:          true
---

Both at work and in my free time I interact with lots of different git repositories - across my machines I usually have about 100 different repositories checked out. 
Maintaining these clones by hand might be possible, but I am way too lazy for that and have written a [go](http://go.dev) program called [ggman](https://github.com/tkw1536/ggman) to manage all of these in a simple way. 

In this post I want to write about how this tool came to be, and give a brief introduction as to it's design and usage.   

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
It was used by setting up a configuration file like this:

```
> Projects 
>> tkw1536
https://github.com/tkw1536/GitManager.git
https://github.com/tkw1536/tkw01536.de.git
https://github.com/tkw1536/guys.wtf.git
```

This file says to create a "Projects" folder.
Then inside it, create a "tkw1536" subfolder. 
Finally clone the three git repositories into this folder. 

I could then use commands like `git-manager pull` or `git-manager push` to pull or push all local repositories. 
This made my job simpler, but I now had to maintain this configuration file. 
I eventually added a function to automatically discover newly cloned repositories and rewrite the configuration file for me. 
This helped a bit more, and I ended up using `git-manager` for all my clones for a couple years. 

## Design goals for ggman

Fast forward again a couple more years to 2019 and I got annoyed at needing to update that configuration file. 
So I decided to redo my repository management again. 

I looked around and there were other tools at the time - however these usually had some downsides:

1. They were limited to one repository provider. 
For example, the [GitHub CLI](http://cli.github.com) only worked with GitHub repositories. 
At the same time, the [GitLab CLI](https://docs.gitlab.com/editor_extensions/gitlab_cli/) only worked with GitLab repositories. 

2. Tools typically encouraged a flat directory structure. 
For example, they cloned repositories directly under a `Projects` folder.
Two repositories that had little to do with each other might end up on disk directly next to each other. 

3. Some tools were only available from within an IDE or GUI.

All of these made felt like annoying downsides.  
As I was starting out with [go](https://go.dev) at the time, I decided to use the opportunity to learn the language properly and start writing my own tool. 
I couldn't think of a good name, I eventually settled on `ggman` - with "man" standing for "manager" and the gs standing for "git" and "go". 

In order to best fit my own workflow, and to prevent me having to rewrite the tool again on the future, I decided on several goals for `ggman`. 
As I have been using `ggman` ever since I designed it without the need for a major rewrite[^1], I consider it successful. 

In particular, I decided `ggman` should:

- be command-line first;
- be simple to install, configure and use;
- encourage an obvious hierarchical directory structure, but remain fully functional with any directory structure;
- remain free of repository provider-specific code; and
- not store any repository-specific data outside of the repositories themselves (enabling the user to switch back to only git at any point).

In order to explain how `ggman` is designed to achieve these goals, I feel like it is best to describe the how to install and use it.
The source code of `ggman` lives [on GitHub](https://github.com/tkw1536/ggman), resulting in a single binary that is dropped into the user's `$PATH` to install.
The binary optionally requires that the user has `git` installed, but will automatically fall back to the [go-git](https://github.com/go-git/go-git) library if not. 
They user can optionally configure several shell aliases, by invoking and evaluating the output of  `ggman shellrc` in their shell's profile.


## Cloning repositories with ggman

Once installed, `ggman` manages all git repositories inside a given root directory, and automatically sets up new repositories relative to the URLs they are cloned from. 
This root folder defaults to `~/Projects`, but can be customized using a `$GGROOT` environment variable.


The first `ggman` command users will likely interact with is one like the following:

```
$ ggman clone https://github.com/tkw1536/ggman.git
Cloning "git@github.com:tkw1536/ggman.git" into "/Users/whoever/github.com/tkw1536/ggman" ...
Cloning into '/Users/whoever/github.com/tkw1536/ggman'...
remote: Enumerating objects: 133, done.
remote: Counting objects: 100% (133/133), done.
remote: Compressing objects: 100% (130/130), done.
remote: Total 133 (delta 16), reused 23 (delta 0), pack-reused 0 (from 0)
Receiving objects: 100% (133/133), 188.95 KiB | 355.00 KiB/s, done.
Resolving deltas: 100% (16/16), done.
```

The `ggman clone` command is intended to clone a repository into the local directory structure. 
It achieves this using several steps:

1.
    Parse the provided into its' so-called URL components.
    Here, the components of a URL are the hostname, the username and '/'-separated elements of the path.
    A username of `git` as well as a trailing suffix of `.git` are dropped.

    Some examples:

    | URL                          | Components                     |
    |------------------------------|--------------------------------|
    | `git@github.com/user/repo`   | `github.com`, `user`, `repo`   |
    | `github.com/hello/world.git` | `github.com`, `hello`, `world` |
    | `github.com/some/repo`       | `gitlab.com`, `some`, `repo`   |
    | `user@server.com:repo.git`   | `server.com`, `user`, `repo`   |
    
    The `ggman comps` command is a utility that allows us to print out components of a specific URL:
    
    ```
    $ ggman comps https://github.com/tkw1536/ggman.git
    github.com
    tkw1536
    ggman 
    ```

    Notice how the components of a URL are identical if cloned via SSH:

    ```
    $ ggman comps git@github.com:tkw1536/ggman.git
    github.com
    tkw1536
    ggman
    ```

    Components are the key abstraction that allow ggman to remain provider independent - as they work with (almost) any repository host. 

2.
    Assign the repository a local path using these components, and create parent folders as needed. 
    In this case the target path would be `$GGROOT/github.com/tkw1536/ggman`.
    The `ggman clone` command above would create `$GGROOT/github.com` and `$GGROOT/github.com/tkw1536` folders as needed.   

3.
    Figure out which URL to clone the repository from.
    This is achieved by turning the components back into a form which git understands. 

    We can inspect this using the `ggman canon` command. 
    In our case:

    ```
    $ ggman canon https://github.com/tkw1536/ggman.git
    git@github.com:tkw1536/ggman.git
    ```

    As you can see, ggman defaults to cloning using an `ssh` clone URL. 
    This can be configured using a so-called `CANSPEC` (short for "canonization specification"), but I won't into detail here.

4.
    Finally invoke the git command to actually clone the repository.

## Finding and performing actions on repositories

But `ggman` can not add clone new repositories. 
It can also perform actions across existing repositories.
Actions in principle take the form ```ggman [FILTERS] ACTION```. 

The supported actions are things which effectively map to plain git commands, such as:

- `ggman ls`, which prints a list of local repositories;
- `ggman pull`, which pulls changes from remotes into the local repositories;
- `ggman push`, which pushes changes to remotes remote repositories.
- `ggman exec COMMAND` which directly invokes an external command.

By default, any action will act on all repositories existing in some sub-directory of `$GGROOT`. 
For example:

```
$ ggman ls
/Users/whoever/github.com/hello/world
/Users/whoever/github.com/tkw1536/ggman
/Users/whoever/github.com/tkw1536/tkw01536.de
/Users/whoever/gitlab.com/lorem/ipsum
```

It is also possible to only act on a subset of repositories using a "FILTER" argument.
The simplest one is the `--for` filter, which fuzzy matches against repositories.

For example:

```
$ ggman --for github.com ls
/Users/whoever/github.com/hello/world
/Users/whoever/github.com/tkw1536/ggman
/Users/whoever/github.com/tkw1536/tkw01536.de
```

This command lists only repositories that match "github.com" in their URL. 
As the matching is fuzzy, it also allows to omit characters or components. 
For example:

```
$ ggman --for lo/ips ls
/Users/whoever/gitlab.com/lorem/ipsum
```

will only match the `lorem/ipsum` repository. 

For convenience ggman provides two shell aliases that make use of the `--for` filter:

-
    `ggcd PATTERN` which finds a repository matching a given pattern and cds into it. 
    For example, `ggcd lo/ip` would cd into `/Users/whoever/gitlab.com/lorem/ipsum`. 

    This makes it extremely simple to find a project belonging to a repository and working on it. 

-
    `ggcode PATTERN`, which is like `ggcd` except that it open a [Visual Studio Code](https://code.visualstudio.com/) instance in the desired directory. 
    This makes it extremely quick to start coding on a specific project, without having to navigate through various user interfaces. 

There are also other filters, but I will omit them in this post. 

## Other ggman functionality

`ggman` has a lot more functionality, but describing everything would make this blog post much longer. 

I would however like to quickly mention a couple of other commands:

-
    `ggman web` which opens the current repository in a web browser. 
    This is useful to quickly use GitHub's web interface to look at issues, or check on the status of CI. 
-
    `ggman relocate` which moves cloned repositories into the paths that `ggman clone` would have cloned them to.
-
    `ggman fix` which updates remote URLs to use their canonical variant. 
 

If you are interested I encourage you to have a look at the [README](https://github.com/tkw1536/ggman?tab=readme-ov-file#ggman) or ask me if you're interested. 

## Conclusion

And that is already all I want to say for now, thank you for reading what basically amounts to a wall of text. 

In summary: I work with lots of git repositories.
I wrote a tool called `ggman` to maintain and expand a local directory structure of all of these. 
It can locate where specific repositories are cloned to, and run operations such as `git pull` or `git push` across all of them. 

Feel free to try it out and give me feedback at https://github.com/tkw1536/ggman. 

[^1]: Unless you count me changing several implementation details under the hood, but I do not.   