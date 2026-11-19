package structcomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestAnalyzerReportsMissingComments(t *testing.T) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "test.go", "package p\n// Good 好的结构体。\ntype Good struct {\n// ID 编号。\nID int\n}\ntype Bad struct { Name string }\n", parser.ParseComments)
	if err != nil { t.Fatal(err) }
	var got []string
	_, err = Analyzer.Run(&analysis.Pass{Fset: fs, Files: []*ast.File{f}, Report: func(d analysis.Diagnostic) { got = append(got, d.Message) }})
	if err != nil { t.Fatal(err) }
	if len(got) != 2 { t.Fatalf("reports = %v", got) }
}
