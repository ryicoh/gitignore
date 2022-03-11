package gitignore

import (
	"fmt"
	"regexp"
	"strings"
)

type Glob struct {
	from        *string
	original    string
	actual      string
	isWhitelist bool
	isOnlyDir   bool
	re          *regexp.Regexp
}

func (g *Glob) hasDoublestarPrefix() bool {
	return strings.HasPrefix(g.actual, "**/") || g.actual == "**"
}

func (g *Glob) Compile() (err error) {
	globStr := g.actual

	reStr := ""
	length := len(globStr)
	for i := 0; i < length; i++ {
		c := globStr[i]
		switch c {
		case '/', '$', '^', '+', '.', '(', ')', '=', '!', '|':
			reStr += "\\" + string(c)
		case '\\':

		case '?':
			reStr += "."

		case '[':
		case ']':
			reStr += string(c)
		case '*':
			var prevChar byte
			if 0 < i && i-1 < len(globStr) {
				prevChar = globStr[i-1]
			}

			var starCount = 1
			for i+1 < len(globStr) && globStr[i+1] == '*' {
				starCount++
				i++
			}

			var nextChar byte
			if i+1 < len(globStr) {
				nextChar = globStr[i+1]
			}

			var isMultiStar = starCount > 1 &&
				(prevChar == 0 || prevChar == '/') &&
				(nextChar == 0 || nextChar == '/')

			if isMultiStar {
				reStr += "((?:[^/]*(?:\\/|$))*)"
				i++
			} else if i == len(globStr)-1 {
				reStr += "((?:[^/]*(?:\\/|$))*)"
			} else {
				reStr += "([^/]*)"
			}

		default:
			reStr += string(c)

		}
	}

	reStr = "^" + reStr + "$"

	g.re, err = regexp.Compile(reStr)
	return err
}

func (g *Glob) match(path string) bool {
	return g.re.MatchString(path)
}

type Gitignore struct {
	root          string
	globs         []Glob
	numIgnores    uint64
	numWhitelists uint64
	matches       []uint
}

func (g *Gitignore) Ignored(path string, isDir bool) bool {
	if len(g.globs) == 0 {
		return false
	}
	return g.isIgnoreStripped(g.strip(path), isDir)
}

func (g *Gitignore) isIgnoreStripped(path string, isDir bool) bool {
	for i := len(g.globs) - 1; i >= 0; i-- {
		glob := g.globs[i]

		if !glob.match(path) {
			continue
		}

		if !glob.isOnlyDir || isDir {
			if glob.isWhitelist {
				return false
			} else {
				return true
			}
		}
		return false
	}

	return false
}

func (g *Gitignore) strip(path string) string {
	// A leading ./ is completely superfluous. We also strip it from
	// our gitignore root path, so we need to strip it from our candidate
	// path too.
	if strings.HasPrefix(path, "./") {
		path = strings.TrimLeft(path, "./")
	}

	// Strip any common prefix between the candidate path and the root
	// of the gitignore, to make sure we get relative matching right.
	// BUT, a file name might not have any directory components to it,
	// in which case, we don't want to accidentally strip any part of the
	// file name.
	//
	// As an additional special case, if the root is just `.`, then we
	// shouldn't try to strip anything, e.g., when path begins with a `.`.
	if g.root != "." && !strings.HasPrefix(path, "/") {
		if strings.HasPrefix(path, g.root) {
			path = strings.TrimLeft(path, g.root)

			if strings.HasPrefix(path, "/") {
				path = strings.TrimLeft(path, "/")
			}
		}
	}
	return path
}

type GitignoreBuilder struct {
	root  string
	globs []Glob
}

func NewGitignoreBuilder(root string) (*GitignoreBuilder, error) {
	return &GitignoreBuilder{
		root:  strings.TrimPrefix(root, "./"),
		globs: []Glob{},
	}, nil
}

func (b *GitignoreBuilder) AddString(from *string, gi string) error {
	for _, line := range strings.Split(gi, "\n") {
		if err := b.AddLine(from, line); err != nil {
			return err
		}
	}
	return nil
}

func (b *GitignoreBuilder) AddLine(from *string, line string) error {
	if strings.HasPrefix(line, "#") {
		return nil
	}

	if !strings.HasSuffix(line, "\\ ") {
		line = strings.TrimSpace(line)
	}

	if line == "" {
		return nil
	}

	glob := Glob{
		from:        from,
		original:    line,
		actual:      "",
		isWhitelist: false,
		isOnlyDir:   false,
	}

	isAbsolute := false
	if strings.HasPrefix(line, "\\!") || strings.HasPrefix(line, "\\#") {
		line = line[1:]
		isAbsolute = line[0] == '/'
	} else {
		if strings.HasPrefix(line, "!") {
			glob.isWhitelist = true
			line = line[1:]
		}
		if strings.HasPrefix(line, "/") {
			// `man gitignore` says that if a glob starts with a slash,
			// then the glob can only match the beginning of a path
			// (relative to the location of gitignore). We achieve this by
			// simply banning wildcards from matching /.
			line = line[1:]
			isAbsolute = true
		}
	}

	// If it ends with a slash, then this should only match directories,
	// but the slash should otherwise not be used while globbing.
	if strings.HasSuffix(line, "/") {
		glob.isOnlyDir = true
		line = line[:len(line)-1]
	}

	glob.actual = line
	// If there is a literal slash, then this is a glob that must match the
	// entire path name. Otherwise, we should let it match anywhere, so use
	// a **/ prefix.
	if !isAbsolute && !strings.Contains(line, "/") {
		// ... but only if we don't already have a **/ prefix.
		if !glob.hasDoublestarPrefix() {
			glob.actual = fmt.Sprintf("**/%s", glob.actual)
		}
	}

	// If the glob ends with `/**`, then we should only match everything
	// inside a directory, but not the directory itself. Standard globs
	// will match the directory. So we add `/*` to force the issue.
	if strings.HasSuffix(glob.actual, "/**") {
		glob.actual = fmt.Sprintf("%s/*", glob.actual)
	}

	if err := glob.Compile(); err != nil {
		return err
	}

	b.globs = append(b.globs, glob)
	return nil
}

func (b *GitignoreBuilder) Build() (*Gitignore, error) {
	var nignore uint64 = 0
	var nwhite uint64 = 0
	for _, g := range b.globs {
		if !g.isWhitelist {
			nignore += 1
		} else {
			nwhite += 1
		}
	}

	return &Gitignore{
		root:          b.root,
		globs:         b.globs,
		numIgnores:    nignore,
		numWhitelists: nwhite,
	}, nil
}
