package mixedreceiver

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mixedreceiver",
	Doc:  "disallow mixing value and pointer receivers on the same struct",
	Run:  runMixedReceivers,
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
