package parser

import (
	"fmt"
	"go/ast"
	"strings"
)

func DetectRequestBodyType(fn *ast.FuncDecl) (map[string]string, error) {
	result := make(map[string]string)

	// Inspect the function body to find ShouldBindJSON calls
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		// Check for a call expression
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if the function being called is ShouldBindJSON
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || selExpr.Sel.Name != "ShouldBindJSON" {
			return true
		}

		// Extract the argument passed to ShouldBindJSON
		if len(callExpr.Args) == 1 {
			// Check if the argument is a unary expression
			unaryExpr, ok := callExpr.Args[0].(*ast.UnaryExpr)
			if ok {
				// Check if the unary expression contains an identifier
				ident, ok := unaryExpr.X.(*ast.Ident)
				if ok {
					// Extract the type of the variable from the declaration
					if ident.Obj != nil && ident.Obj.Decl != nil {
						valueSpec, ok := ident.Obj.Decl.(*ast.ValueSpec)
						if ok && len(valueSpec.Type.(*ast.Ident).Name) > 0 {
							typeName := valueSpec.Type.(*ast.Ident).Name
							result[typeName] = "" // struct name
						}
					}
				}
			}
		}

		return true
	})

	if len(result) > 0 {
		return result, nil
	}
	return nil, fmt.Errorf("no ShouldBindJSON found")
}

func DetectResponseModel(fn *ast.FuncDecl) map[string]string {
	responses := make(map[string]string)

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) < 2 {
			return true
		}

		// Check for selector expression: c.JSON
		if selExpr, ok := call.Fun.(*ast.SelectorExpr); ok && selExpr.Sel.Name == "JSON" {
			// Attempt to get status code and response model type
			statusCodeLit, ok1 := call.Args[0].(*ast.BasicLit)
			responseExpr := call.Args[1]

			if ok1 {
				status := strings.Trim(statusCodeLit.Value, "\"")

				switch expr := responseExpr.(type) {
				case *ast.CompositeLit:
					if ident, ok := expr.Type.(*ast.Ident); ok {
						responses[status] = ident.Name
					}
				case *ast.Ident:
					responses[status] = expr.Name
				}
			}
		}

		return true
	})

	return responses
}
