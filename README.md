# Git Exec

A library that wraps git commands and returns structured output.

## Example Usage

Install using standard go get:

```go
go get github.com/instruqt/git-exec
```

Use in your project:

```go
package main

import (
  ge "github.com/instruqt/git-exec"
)

func main() {
  git, err := ge.New()
  if err != nil {}

  message, err := git.Init()
  if err != nil {}

  fmt.Println(message)

  err = os.WriteFile("file.txt", []byte("Hello, World!"), 0644)
  if err != nil {}

  err = git.Add("file.txt")
  if err != nil {}

  files, err := git.Status()
  if err != nil {}

  fmt.Printf("%#v", files)

  err = git.Commit("Initial commit", "Name", "name@gmail.com")
  if err != nil {}
}
```