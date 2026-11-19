package funccomment

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestAnalyzerReportsMissingAndInvalidComments(t *testing.T) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "test.go", `package p
func Good() {}
// Bad only English.
func Bad() {}
// Valid 有效说明。
func Valid() {}
`, parser.ParseComments)
	if err != nil { t.Fatal(err) }
	var messages []string
	_, err = Analyzer.Run(&analysis.Pass{Fset: fileSet, Files: []*ast.File{file}, Report: func(d analysis.Diagnostic) { messages = append(messages, d.Message) }})
	if err != nil { t.Fatal(err) }
	if len(messages) != 2 { t.Fatalf("reported %d issues, want 2: %v", len(messages), messages) }
}
