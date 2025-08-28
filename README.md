# go-starter

A minimal Go starter repository with VSCode configuration, linting, testing, a Makefile for custom commands, and pre-push hooks.

## Prerequisites

- Go 1.20+
- [golangci-lint](https://golangci-lint.run)
- Make
- VSCode with the Go extension installed (pre-release version is required to support golangci-lint v2)

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/iluhanec/go-starter.git
   cd go-starter
   ```

2. Configure Git hooks:

   ```sh
   git config core.hooksPath .githooks
   ```

3. Ensure the pre-push hook is executable:

   ```sh
   chmod +x .githooks/pre-push
   ```

4. Open in VSCode:
   ```sh
   code .
   ```

VSCode will automatically format and organize imports on save.

## Running & Custom Commands

This repo uses a **Makefile** to define common tasks, similar to `package.json` scripts in Node.js.

| Command      | Description         |
| ------------ | ------------------- |
| `make run`   | Run the application |
| `make build` | Compile the binary  |
| `make test`  | Execute all tests   |
| `make lint`  | Run code linter     |

Example:

```sh
make run
```

## Building a Binary

You can compile your application into a standalone executable:

```sh
make build
./go-starter
```

Or directly with Go:

```sh
go build -o go-starter main.go
./go-starter

```
