# Committer

A simple CLI to generate git commit messages with AI

## Installation

```sh
brew tap ethn1ee/committer
brew install committer
```

## Usage

```sh
committer gen
```

To generate commit message and sync all at once,

```sh
git add .
git commit -m $(committer gen)
git push
```
