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
		{name: "English suffix with Chinese on next line", src: "// Foo xxxxj\n// 执行处理。\nfunc Foo() {}", want: true},
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
			if got := hasNamePrefixedComment(fn.Doc, "Foo"); got != tt.want {
				t.Fatalf("hasNamePrefixedComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunChecksConstants(t *testing.T) {
	src := `package p

// GoodConst 有中文说明。
const GoodConst = 1

const BadConst = 2

const (
	// GroupedConst 分组常量说明。
	GroupedConst = 3
	MissingGroupedConst = 4
)
`
	file, err := parser.ParseFile(token.NewFileSet(), "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok.String() == "const" {
			for _, spec := range gen.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				for _, name := range valueSpec.Names {
					if name.Name == "GoodConst" || name.Name == "GroupedConst" {
						doc := valueSpec.Doc
						if doc == nil && len(gen.Specs) == 1 {
							doc = gen.Doc
						}
						if !hasNamePrefixedComment(doc, name.Name) {
							t.Errorf("expected %s to have a valid comment", name.Name)
						}
					}
				}
			}
		}
	}
}
