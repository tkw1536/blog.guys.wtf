---
title:          Why I wrote a tool to manage all my repositories 
date:           2025-07-13
author:         Tom Wiesing 
authorLink:     https://tkw01536.de

description:    An overview of how I designed and wrote a tool called ggman to manage all my git repositories. 

draft:          true
---

This is still a draft, copy-pasted from the ggman readme. 
Need to rework this.


A golang tool that can manage all your git repositories. 

## What Is ggman?

When you only have a couple of git repositories that you work on it is perfectly feasible to manage them by using `git clone`, `git pull` and friends. 
However once the number of repositories grows beyond a small number this can become tedious:

- It is hard to find which folder a repository has been cloned to
- Getting an overview of what is cloned and what is not is hard
- It's not easily possible to perform actions on more than one repo at once, e.g. `git pull`

This is the problem `ggman` is designed to solve. 
It allows one to:

- Maintain and expand a local directory structure of multiple repositories
- Run actions (such as `git clone`, `git pull`) on groups of repositories

## Why ggman?

While similar tools exist these commonly have a lot of downsides:

- they enforce a flat directory structure;
- they are limited to one repository provider (such as GitHub or GitLab); or
- they are only available from within an IDE or GUI.

ggman considers these as major downsides. 
The goals and principles of ggman are:

- to be command-line first;
- to be simple to install, configure and use;
- to encourage an obvious hierarchical directory structure, but remain fully functional with any directory structure;
- to remain free of provider-specific code; and
- to not store any repository-specific data outside of the repositories themselves (enabling the user to switch back to only git at any point).
