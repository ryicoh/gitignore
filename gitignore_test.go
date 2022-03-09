package gitignore

import "testing"

func giFromStr(root string, s string) (*Gitignore, error) {
	builder, err := NewGitignoreBuilder(root)
	if err != nil {
		return nil, err
	}
	if err := builder.AddString(nil, s); err != nil {
		return nil, err
	}
	return builder.Build()
}

const ROOT = "/home/foobar/rust/rg"

func TestIgnore(t *testing.T) {
	testCases := []struct {
		desc     string
		root     string
		gi       string
		path     string
		isDir    bool
		expected bool
	}{
		{"ig1", ROOT, "months", "months", false, true},
		{"ig2", ROOT, "*.lock", "Cargo.lock", false, true},
		{"ig3", ROOT, "*.rs", "src/main.rs", false, true},
		{"ig4", ROOT, "src/*.rs", "src/main.rs", false, true},
		{"ig5", ROOT, "/*.c", "cat-file.c", false, true},
		{"ig6", ROOT, "/src/*.rs", "src/main.rs", false, true},
		{"ig7", ROOT, "!src/main.rs\n*.rs", "src/main.rs", false, true},
		{"ig8", ROOT, "foo/", "foo", true, true},
		{"ig9", ROOT, "**/foo", "foo", false, true},
		{"ig10", ROOT, "**/foo", "src/foo", false, true},
		{"ig11", ROOT, "**/foo/**", "src/foo/bar", false, true},
		{"ig12", ROOT, "**/foo/**", "wat/src/foo/bar/baz", false, true},
		{"ig13", ROOT, "**/foo/bar", "foo/bar", false, true},
		{"ig14", ROOT, "**/foo/bar", "src/foo/bar", false, true},
		{"ig15", ROOT, "abc/**", "abc/x", false, true},
		{"ig16", ROOT, "abc/**", "abc/x/y", false, true},
		{"ig17", ROOT, "abc/**", "abc/x/y/z", false, true},
		{"ig18", ROOT, "a/**/b", "a/b", false, true},
		{"ig19", ROOT, "a/**/b", "a/x/b", false, true},
		{"ig20", ROOT, "a/**/b", "a/x/y/b", false, true},
		{"ig21", ROOT, "\\!xy", "!xy", false, true},
		{"ig22", ROOT, "\\#foo", "#foo", false, true},
		{"ig23", ROOT, "foo", "./foo", false, true},
		{"ig24", ROOT, "target", "grep/target", false, true},
		{"ig25", ROOT, "Cargo.lock", "./tabwriter-bin/Cargo.lock", false, true},
		{"ig26", ROOT, "/foo/bar/baz", "./foo/bar/baz", false, true},
		{"ig27", ROOT, "foo/", "xyz/foo", true, true},
		{"ig28", "./src", "/llvm/", "./src/llvm", true, true},
		{"ig29", ROOT, "node_modules/ ", "node_modules", true, true},
		{"ig30", ROOT, "**/", "foo/bar", true, true},
		{"ig31", ROOT, "path1/*", "path1/foo", false, true},
		{"ig32", ROOT, ".a/b", ".a/b", false, true},
		{"ig33", "./", ".a/b", ".a/b", false, true},
		{"ig34", ".", ".a/b", ".a/b", false, true},
		{"ig35", "./.", ".a/b", ".a/b", false, true},
		{"ig36", "././", ".a/b", ".a/b", false, true},
		{"ig37", "././.", ".a/b", ".a/b", false, true},
		{"ig38", ROOT, "\\[", "[", false, true},
		{"ig39", ROOT, "\\?", "?", false, true},
		{"ig40", ROOT, "\\*", "*", false, true},
		{"ig41", ROOT, "\\a", "a", false, true},
		{"ig42", ROOT, "s*.rs", "sfoo.rs", false, true},
		{"ig43", ROOT, "**", "foo.rs", false, true},
		{"ig44", ROOT, "**/**/*", "a/foo.rs", false, true},
		{"ignot1", ROOT, "amonths", "months", false, false},
		{"ignot2", ROOT, "monthsa", "months", false, false},
		{"ignot3", ROOT, "/src/*.rs", "src/grep/src/main.rs", false, false},
		{"ignot4", ROOT, "/*.c", "mozilla-sha1/sha1.c", false, false},
		{"ignot5", ROOT, "/src/*.rs", "src/grep/src/main.rs", false, false},
		{"ignot6", ROOT, "*.rs\n!src/main.rs", "src/main.rs", false, false},
		{"ignot7", ROOT, "foo/", "foo", false, false},
		{"ignot8", ROOT, "**/foo/**", "wat/src/afoo/bar/baz", false, false},
		{"ignot9", ROOT, "**/foo/**", "wat/src/fooa/bar/baz", false, false},
		{"ignot10", ROOT, "**/foo/bar", "foo/src/bar", false, false},
		{"ignot11", ROOT, "#foo", "#foo", false, false},
		{"ignot12", ROOT, "\n\n\n", "foo", false, false},
		{"ignot13", ROOT, "foo/**", "foo", true, false},
		{"ignot14", "./third_party/protobuf", "m4/ltoptions.m4", "./third_party/protobuf/csharp/src/packages/repositories.config", false, false},
		{"ignot15", ROOT, "!/bar", "foo/bar", false, false},
		{"ignot16", ROOT, "*\n!**/", "foo", true, false},
		{"ignot17", ROOT, "src/*.rs", "src/grep/src/main.rs", false, false},
		{"ignot18", ROOT, "path1/*", "path2/path1/foo", false, false},
		{"ignot19", ROOT, "s*.rs", "src/foo.rs", false, false},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			gi, err := giFromStr(tC.root, tC.gi)
			if err != nil {
				t.Fatal(err)
			}
			actual := gi.IsIgnore(tC.path, tC.isDir)
			if actual != tC.expected {
				t.Fatalf("%s: root(%s), gi(%s), path(%s) actual(%v), expect(%v)\n",
					tC.desc, tC.root, tC.gi, tC.path, actual, tC.expected)
			}
		})
	}
}
