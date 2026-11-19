package funccomment

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "funccomment",
	Doc:  "require every function to have a name-prefixed Chinese comment",
	Run:  run,
}

var ConstAnalyzer = &analysis.Analyzer{
	Name: "constcomment",
	Doc:  "require every constant to have a name-prefixed Chinese comment",
	Run:  runConstants,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			switch node := decl.(type) {
			case *ast.FuncDecl:
				if node.Name.Name == "init" {
					continue
				}
				checkFunctionComment(pass, node)
			}
		}
	}

	return nil, nil
}

func runConstants(pass *analysis.Pass) (any, error) {
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

func checkFunctionComment(pass *analysis.Pass, fn *ast.FuncDecl) {
	name := fn.Name.Name
	if fn.Doc == nil || len(fn.Doc.List) == 0 {
		pass.Reportf(fn.Pos(), "function %s must have a comment starting with its name and containing Chinese text", name)
		return
	}
	if !hasNamePrefixedComment(fn.Doc, name) {
		pass.Reportf(fn.Doc.Pos(), "comment for %s must start with %s and contain Chinese text", name, name)
	}
}

func checkConstantComment(pass *analysis.Pass, pos token.Pos, doc *ast.CommentGroup, name string) {
	if doc == nil || len(doc.List) == 0 {
		pass.Reportf(pos, "constant %s must have a comment starting with its name and containing Chinese text", name)
		return
	}
	if !hasNamePrefixedComment(doc, name) {
		pass.Reportf(doc.Pos(), "comment for constant %s must start with %s and contain Chinese text", name, name)
	}
}

func hasNamePrefixedComment(doc *ast.CommentGroup, name string) bool {
	if doc == nil || len(doc.List) == 0 {
		return false
	}

	comment := strings.TrimSpace(doc.Text())
	firstLine := comment
	if lineEnd := strings.IndexByte(comment, '\n'); lineEnd >= 0 {
		firstLine = comment[:lineEnd]
	}

	if !strings.HasPrefix(firstLine, name) {
		return false
	}

	rawRemainder := strings.TrimPrefix(firstLine, name)
	if rawRemainder != "" {
		firstRemainder, _ := utf8.DecodeRuneInString(rawRemainder)
		if !unicode.IsSpace(firstRemainder) && firstRemainder != ':' && firstRemainder != '：' {
			return false
		}
	}

	remainder := strings.TrimSpace(strings.TrimPrefix(comment, name))
	remainder = strings.TrimSpace(strings.TrimPrefix(remainder, ":"))
	remainder = strings.TrimSpace(strings.TrimPrefix(remainder, "："))
	if remainder == "" {
		return false
	}

	return strings.ContainsFunc(remainder, func(r rune) bool {
		return unicode.Is(unicode.Han, r)
	})
}
