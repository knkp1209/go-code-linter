package funccomment

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "funccomment",
	Doc:  "require every function to have a name-prefixed comment",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			name := fn.Name.Name

			// init 函数不检查
			if name == "init" {
				continue
			}

			if fn.Doc == nil || len(fn.Doc.List) == 0 {
				pass.Reportf(
					fn.Pos(),
					"function %s must have a comment starting with its name",
					name,
				)
				continue
			}

			comment := strings.TrimSpace(
				strings.TrimPrefix(fn.Doc.List[0].Text, "//"),
			)

			// 支持 "// FunctionName ..." 和 "// FunctionName"
			if comment != name &&
				!strings.HasPrefix(comment, name+" ") &&
				!strings.HasPrefix(comment, name+":") {
				pass.Reportf(
					fn.Doc.Pos(),
					"comment for %s must start with %s",
					name,
					name,
				)
			}
		}
	}

	return nil, nil
}
