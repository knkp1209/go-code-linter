package funccomment

import (
	"go/ast"

	"github.com/knkp1209/go-code-linter/internal/commentcheck"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "funccomment",
	Doc:  "require every function to have a name-prefixed Chinese comment",
	Run:  runFunctionComments,
}

func runFunctionComments(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name == "init" {
				continue
			}
			checkFunctionComment(pass, fn)
		}
	}
	return nil, nil
}

func checkFunctionComment(pass *analysis.Pass, fn *ast.FuncDecl) {
	name := fn.Name.Name
	if fn.Doc == nil || len(fn.Doc.List) == 0 {
		pass.Reportf(fn.Pos(), "function %s must have a comment starting with its name and containing Chinese text", name)
		return
	}
	if !commentcheck.HasNamePrefixedComment(fn.Doc, name) {
		pass.Reportf(fn.Doc.Pos(), "comment for %s must start with %s and contain Chinese text", name, name)
	}
}
