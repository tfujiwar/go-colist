# go-colist
`colist` is a command to show GitHub code owners of files that were modified in a current branch.

## Example
```console
$ cat .github/CODEOWNERS
*.cc @org/cc-reviewer-team
*.go @org/go-reviewer-team
*.rs @org/rs-reviewer-team

$ git diff main... --name-only
foo.go
bar.go

$ colist
*.go : @org/go-developer-team
```

## Install
```
go install github.com/tfujiwar/go-colist/cmd/colist@latest
```

## Usage
```
Usage:
  colist [flags]                        : compare with remote or local main branch
  colist [flags] <BASE_BRANCH>          : compare with remote or local <BASE_BRANCH>
  colist [flags] <REMOTE> <BASE_BRANCH> : compare with <REMOTE>/<BASE_BRANCH>

Flags:
  -o, --output text|json : output format
  -d, --dir <DIR>        : repository directory
  -v, --verbose          : show debug log
  -h, --help             : show this message
```
