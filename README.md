# gitignore
[![Go](https://github.com/ryicoh/gitignore/actions/workflows/go.yml/badge.svg)](https://github.com/ryicoh/gitignore/actions/workflows/go.yml)

.gitignore pattern matching in Go.


# Installation
```go

import "github.com/ryicoh/gitignore"

```

# Example

```go
package main

import (
	"fmt"

	"github.com/ryicoh/gitignore"
)

func main() {
	gi, _ := gitignore.NewGitignoreFromDir(".")
	ignored := gi.Ignored("cmd/main.go", false)
	fmt.Println(ignored)
}
```

# Acknowledgments

Most of the implementation was based on ripgrep's [gitignore.rs](https://github.com/BurntSushi/ripgrep/blob/master/crates/ignore/src/gitignore.rs)
