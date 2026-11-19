package mixedreceiver

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestAnalyzerReportsMixedReceivers(t *testing.T) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "test.go", "package p\ntype Snapshot struct{}\nfunc (Snapshot) Value(){}\nfunc (*Snapshot) Pointer(){}\n", parser.ParseComments)
	if err != nil { t.Fatal(err) }
	var got []string
	_, err = Analyzer.Run(&analysis.Pass{Fset: fs, Files: []*ast.File{f}, Report: func(d analysis.Diagnostic) { got = append(got, d.Message) }})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 { t.Fatalf("reports = %v", got) }
}
