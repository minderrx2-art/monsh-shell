# Monty Shell

A small Unix-like shell written in Go.

Monty Shell is a toy project built to learn how shells work under the hood. The goal to understand the core components that make a shell work.

## Features (so far)

- REPL
- Execute binaries
- Built-in commands
  - `cd`
  - `pwd`
  - `echo`
  - `type`
  - `exit`
- `PATH` executable lookup
- Redirection (`>` `>>` `<`) also (`1> 2> 2>> 0<`)

## Running

```bash
go run ./cmd/shell
```

## Example

```bash
$ pwd
/home/monty
$ echo hello world
hello world
$ ls
README.md
go.mod
internal
$ echo "hello" > file.txt
$ cat file.txt
hello
$ echo "world" >> file.txt
$ cat file.txt
hello world
$ type ls
ls is /usr/bin/ls
$ cd ..
```

## Architecture
1. **Lexer** converts user input into tokens.
2. **Parser** builds an abstract syntax tree (AST).
3. **Executor** walks the AST.
4. **Built-ins** are handled internally.
5. **External commands** are located using `PATH` and executed as child processes.

```
stdin -> Lexer -> Parser -> AST -> Executor
                                      ├── Built-ins
                                      └── External Process
```