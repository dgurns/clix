# CLIX

An LLM-powered CLI agent built with Go.

## Demo

<video src="https://github.com/dgurns/clix/assets/1173791/96f83fdf-5eed-48fb-8532-c3af242a0659" width="100%"></video>

## Dependencies

- [Go](https://golang.org/doc/install) `v1.21`
- [OpenAI API Key](https://platform.openai.com)

## Installation

After cloning the repo to your computer, run `make install`. This will build the `clix` binary and install it to your `$GOPATH/bin` directory.

## Usage

In your terminal, simply run `clix`. You'll be prompted to enter your OpenAI API key. Then you can start asking `clix` to help you with tasks on your computer.

Currently, `clix` uses `gpt-4o` for LLM.

## Reset

If you want to clear the current API key and remove the `clix` binary, run `make reset`.

## Roadmap

- [ ] Support more models like Claude and Llama
- [ ] Enable editing commands that the LLM suggests
- [ ] Pass an initial command like `clix "How do I x?"`
- [ ] Stream command output to Clix stdout
