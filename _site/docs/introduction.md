---
slug: /
---

# Introduction

`mani` is a CLI tool that helps you manage multiple repositories. It's useful when you are working with microservices, multi-project systems, many libraries or just a bunch of repositories and want a central place for pulling all repositories and running commands over them.

## Features

- Clone multiple repositories in one command
- Declarative configuration
- Run custom or ad-hoc commands over multiple repositories
- Flexible filtering
- Portable, no dependencies
- Supports auto-completion

## Example

You specify repository and commands in a config file:

```yaml title="mani.yaml"
projects:
  pinto:
    url: https://github.com/alajmo/pinto.git
    desc: A vim theme editor
    tags: [frontend, node]

  dashgrid:
    url: https://github.com/alajmo/dashgrid.git
    desc: A highly customizable drag-and-drop grid
    tags: [lib, node]

tasks:
  git-status:
    desc: show working tree status
    cmd: git status

  git-create:
    desc: create branch
    output: text
    env:
      branch: main
    cmd: git checkout -b $branch
```

run `mani sync` to clone the repositories:

```bash
$ mani sync
✓ pinto
✓ dashgrid

All projects synced
```

and then run the commands over all or a subset of the repositories:

```bash
# Target repositories that have the tag 'node'
$ mani run git-status --tags node

┌──────────┬─────────────────────────────────────────────────┐
│ Project  │ git-status                                      │
├──────────┼─────────────────────────────────────────────────┤
│ pinto    │ On branch master                                │
│          │ Your branch is up to date with 'origin/master'. │
│          │                                                 │
│          │ nothing to commit, working tree clean           │
├──────────┼─────────────────────────────────────────────────┤
│ dashgrid │ On branch master                                │
│          │ Your branch is up to date with 'origin/master'. │
│          │                                                 │
│          │ nothing to commit, working tree clean           │
└──────────┴─────────────────────────────────────────────────┘



# Target project 'pinto'
$ mani run git-create --projects pinto branch=dev

[pinto] TASK [git-create: create branch] *******************

Switched to a new branch 'dev'
```
