# Committer

A simple CLI to generate git commit messages with AI

![Demo](./assets/demo.gif)

## Installation

```sh
brew tap ethn1ee/committer
brew install committer
```

## GEMINI API Key

You need to set the `GEMINI_API_KEY` environment variable to use this tool.
You can get an API key from the [Google Cloud Console](https://console.cloud.google.com/apis/credentials).

```sh
export GEMINI_API_KEY="your_api_key_here"
```

## Usage

Generate a commit message and print:

```sh
committer gen
```

Generate a commit message and commit:

```sh
committer gen -c
```

Generate a commit message and commit + push

```sh
committer gen -p
```
