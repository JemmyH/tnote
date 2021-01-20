# tnote
`tnote` is a terminal notebook app with password, in which you can record some notes in your terminal.

## Install
```bash
$ cd $GOPATH && mkdir jemmyh && cd jemmyh && git clone https://github.com/JemmyH/tnote.git
$ cd tnote && go build -o tnote && mv tnote $GOPATH/bin/
```

or install with `go get`:
```bash
go get -u -v github.com/JemmyH/tnote
```

## Usage
```bash
$ tnote
Terminal Notebook is CLI App, which is implemented by Golang.

Usage:
   [command]

Available Commands:
  add         Add a note to your notebook
  create      Create a notebook
  delete      Delete notes prefixed with id
  help        Help about any command
  print       Print notes in your notebook
  version     Show version

Flags:
  -h, --help           help for this command
  -o, --owner string   owner of notebook
```

For each subcommand, use `tnote help xxx` for detail usage.

#### create a notebook
```bash
tnote create --owner=xxx
```

#### add a note
```bash
tnote add --owner=xxx --content=hello world
```

#### delete notes which has prefixed id
```bash
tnote delete --owner=xxx --id=20200102
```

### print notes from past to now
```bash
tnote --owner=xxx print -v
```