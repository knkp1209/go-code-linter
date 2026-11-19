package commentcheck

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestHasNamePrefixedComment(t *testing.T) {
	tests := []struct {
		name, src string
		want      bool
	}{
		{"missing", "func Foo() {}", false},
		{"missing description", "// Foo\nfunc Foo() {}", false},
		{"Chinese description", "// Foo 执行处理。\nfunc Foo() {}", true},
		{"Chinese on next line", "// Foo xxxxj\n// 执行处理。\nfunc Foo() {}", true},
		{"wrong name", "// Bar 执行处理。\nfunc Foo() {}", false},
		{"block comment", "/* Foo 执行处理。 */\nfunc Foo() {}", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseFile(token.NewFileSet(), "test.go", "package p\n\n"+tt.src, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			fn := file.Decls[0].(*ast.FuncDecl)
			if got := HasNamePrefixedComment(fn.Doc, "Foo"); got != tt.want {
				t.Fatalf("HasNamePrefixedComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
