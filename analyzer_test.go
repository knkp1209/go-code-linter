package funccomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestHasNamePrefixedComment(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want bool
	}{
		{name: "missing", src: "func Foo() {}", want: false},
		{name: "missing description", src: "// Foo\nfunc Foo() {}", want: false},
		{name: "Chinese description", src: "// Foo 执行处理。\nfunc Foo() {}", want: true},
		{name: "Chinese description on next line", src: "// Foo\n// 执行处理。\nfunc Foo() {}", want: true},
		{name: "colon", src: "// Foo: 执行处理。\nfunc Foo() {}", want: true},
		{name: "full width colon", src: "// Foo：执行处理。\nfunc Foo() {}", want: true},
		{name: "English description", src: "// Foo does something.\nfunc Foo() {}", want: false},
		{name: "wrong name", src: "// Bar 执行处理。\nfunc Foo() {}", want: false},
		{name: "block comment", src: "/* Foo 执行处理。 */\nfunc Foo() {}", want: true},
		{name: "multiline block comment", src: "/*\nFoo 执行处理。\n*/\nfunc Foo() {}", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseFile(token.NewFileSet(), "test.go", "package p\n\n"+tt.src, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			fn := file.Decls[0].(*ast.FuncDecl)
			if got := hasNamePrefixedComment(fn, "Foo"); got != tt.want {
				t.Fatalf("hasNamePrefixedComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
