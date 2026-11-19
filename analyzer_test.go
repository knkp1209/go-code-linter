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

func TestStructComments(t *testing.T) {
	src := `package p

// User 用户信息。
type User struct {
	// ID 用户编号。
	ID int
	// Name 用户名称。
	Name string
}

type MissingDoc struct {
	// ID 编号。
	ID int
	Name string
}
`
	file, err := parser.ParseFile(token.NewFileSet(), "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		typeSpec := gen.Specs[0].(*ast.TypeSpec)
		structType := typeSpec.Type.(*ast.StructType)
		typeDoc := typeSpec.Doc
		if typeDoc == nil && len(gen.Specs) == 1 {
			typeDoc = gen.Doc
		}
		if typeSpec.Name.Name == "User" && !hasNamePrefixedComment(typeDoc, "User") {
			t.Error("expected User to have a valid comment")
		}
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				got := hasNamePrefixedComment(field.Doc, name.Name)
				want := name.Name != "Name" || typeSpec.Name.Name == "User"
				if got != want {
					t.Errorf("field %s comment valid = %v, want %v", name.Name, got, want)
				}
			}
		}
	}
}
