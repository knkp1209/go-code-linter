package funccomment

import (
	"go/ast"
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
					"function %s must have a comment starting with its name and containing Chinese text",
					name,
				)
				continue
			}

			if !hasNamePrefixedComment(fn, name) {
				pass.Reportf(
					fn.Doc.Pos(),
					"comment for %s must start with %s and contain Chinese text",
					name,
					name,
				)
			}
		}
	}

	return nil, nil
}

func hasNamePrefixedComment(fn *ast.FuncDecl, name string) bool {
	if fn.Doc == nil || len(fn.Doc.List) == 0 {
		return false
	}

	comment := strings.TrimSpace(fn.Doc.Text())
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
