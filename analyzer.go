package codelint

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var FunctionCommentAnalyzer = &analysis.Analyzer{
	Name: "funccomment",
	Doc:  "require every function to have a name-prefixed Chinese comment",
	Run:  run,
}

var ConstantCommentAnalyzer = &analysis.Analyzer{
	Name: "constcomment",
	Doc:  "require every constant to have a name-prefixed Chinese comment",
	Run:  runConstants,
}

var StructCommentAnalyzer = &analysis.Analyzer{
	Name: "structcomment",
	Doc:  "require every struct and field to have a name-prefixed Chinese comment",
	Run:  runStructs,
}

var MixedReceiverAnalyzer = &analysis.Analyzer{
	Name: "mixedreceiver",
	Doc:  "disallow mixing value and pointer receivers on the same struct",
	Run:  runMixedReceivers,
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

func runStructs(pass *analysis.Pass) (any, error) {
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

func runMixedReceivers(pass *analysis.Pass) (any, error) {
	type receiverKinds struct {
		value, pointer bool
		pos            token.Pos
	}

	structs := make(map[string]token.Pos)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				if _, ok := typeSpec.Type.(*ast.StructType); ok {
					structs[typeSpec.Name.Name] = typeSpec.Name.Pos()
				}
			}
		}
	}

	kinds := make(map[string]receiverKinds)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			name, pointer, ok := receiverTypeName(fn.Recv.List[0].Type)
			if !ok {
				continue
			}
			pos, ok := structs[name]
			if !ok {
				continue
			}
			kind := kinds[name]
			kind.pos = pos
			kind.pointer = kind.pointer || pointer
			kind.value = kind.value || !pointer
			kinds[name] = kind
		}
	}

	for name, kind := range kinds {
		if kind.value && kind.pointer {
			pass.Reportf(kind.pos, "struct %s has both value and pointer receiver methods; use one receiver style consistently", name)
		}
	}
	return nil, nil
}

func receiverTypeName(expr ast.Expr) (name string, pointer bool, ok bool) {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name, false, true
	case *ast.StarExpr:
		name, _, ok := receiverTypeName(expr.X)
		return name, true, ok
	case *ast.ParenExpr:
		return receiverTypeName(expr.X)
	case *ast.IndexExpr:
		return receiverTypeName(expr.X)
	case *ast.IndexListExpr:
		return receiverTypeName(expr.X)
	default:
		return "", false, false
	}
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

func checkStructComment(pass *analysis.Pass, pos token.Pos, doc *ast.CommentGroup, name string) {
	if doc == nil || len(doc.List) == 0 {
		pass.Reportf(pos, "struct or field %s must have a comment starting with its name and containing Chinese text", name)
		return
	}
	if !hasNamePrefixedComment(doc, name) {
		pass.Reportf(doc.Pos(), "comment for struct or field %s must start with %s and contain Chinese text", name, name)
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
