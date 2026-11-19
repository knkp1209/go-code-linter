package constcomment

import (
	"go/ast"
	"go/token"

	"github.com/knkp1209/go-code-linter/internal/commentcheck"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "constcomment",
	Doc:  "require every constant to have a name-prefixed Chinese comment",
	Run:  runConstantComments,
}

func runConstantComments(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			for _, spec := range gen.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				for _, name := range valueSpec.Names {
					doc := valueSpec.Doc
					if doc == nil && len(gen.Specs) == 1 {
						doc = gen.Doc
					}
					checkConstantComment(pass, name.Pos(), doc, name.Name)
				}
			}
		}
	}
	return nil, nil
}

func checkConstantComment(pass *analysis.Pass, pos token.Pos, doc *ast.CommentGroup, name string) {
	if doc == nil || len(doc.List) == 0 {
		pass.Reportf(pos, "constant %s must have a comment starting with its name and containing Chinese text", name)
		return
	}
	if !commentcheck.HasNamePrefixedComment(doc, name) {
		pass.Reportf(doc.Pos(), "comment for constant %s must start with %s and contain Chinese text", name, name)
	}
}
