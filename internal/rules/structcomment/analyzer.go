package structcomment

import (
	"go/ast"
	"go/token"

	"github.com/knkp1209/go-code-linter/internal/commentcheck"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "structcomment",
	Doc:  "require every struct and field to have a name-prefixed Chinese comment",
	Run:  runStructComments,
}

func runStructComments(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				doc := typeSpec.Doc
				if doc == nil && len(gen.Specs) == 1 {
					doc = gen.Doc
				}
				checkStructComment(pass, typeSpec.Name.Pos(), doc, typeSpec.Name.Name)
				for _, field := range structType.Fields.List {
					for _, name := range field.Names {
						checkStructComment(pass, name.Pos(), field.Doc, name.Name)
					}
				}
			}
		}
	}
	return nil, nil
}

func checkStructComment(pass *analysis.Pass, pos token.Pos, doc *ast.CommentGroup, name string) {
	if doc == nil || len(doc.List) == 0 {
		pass.Reportf(pos, "struct or field %s must have a comment starting with its name and containing Chinese text", name)
		return
	}
	if !commentcheck.HasNamePrefixedComment(doc, name) {
		pass.Reportf(doc.Pos(), "comment for struct or field %s must start with %s and contain Chinese text", name, name)
	}
}
