package constcomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestAnalyzerReportsMissingComment(t *testing.T) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "test.go", "package p\n// Good 有效说明。\nconst Good = 1\nconst Bad = 2\n", parser.ParseComments)
	if err != nil { t.Fatal(err) }
	var got []string
	_, err = Analyzer.Run(&analysis.Pass{Fset: fs, Files: []*ast.File{f}, Report: func(d analysis.Diagnostic) { got = append(got, d.Message) }})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 { t.Fatalf("reports = %v", got) }
}
