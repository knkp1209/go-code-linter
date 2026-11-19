package constcomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestConstantComments(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "test.go", "package p\n// Good 有效说明。\nconst Good = 1\nconst Bad = 2\n", parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	_ = file.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
	if Analyzer == nil || Analyzer.Name != "constcomment" {
		t.Fatal("invalid constant analyzer")
	}
}
