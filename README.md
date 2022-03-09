# gitignore

.gitignore pattern matching in Go.


# Installation
```go

import "github.com/ryicoh/gitignore"

```

# Example

```
package main

import (
	"fmt"

	"github.com/ryicoh/gitignore"
)

func main() {
	var err error
	builder, err := gitignore.NewGitignoreBuilder(".")
	err := builder.AddString(nil, "*.go")
	gi, err := builder.Build()
	ignored := gi.Ignored("cmd/main.go", false)
	fmt.Println(ignored)
}
```

# Acknowledgments

Most of the implementation was based on ripgrep's [gitignore.rs](https://github.com/BurntSushi/ripgrep/blob/master/crates/ignore/src/gitignore.rs)
