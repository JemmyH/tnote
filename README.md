# terminal_note
`terminal_note` is a terminal notebook app, in which you can record some notes in your terminal and set password for it.

## Install
```go
$ cd $GOPATH && git clone https://github.com/JemmyH/terminal_note.git
$ cd terminal_note && go build -o tnote && mv tnote $GOPATH/bin/
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
tnote add --owner=xxx --content=我的头发长天下我为王
```

#### delete notes which has prefixed id
```bash
tnote delete --owner=xxx --id=20200102
```

### print notes from past to now
```bash
tnote --owner=xxx print -v
```