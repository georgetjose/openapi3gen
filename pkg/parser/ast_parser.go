package parser

import (
	"fmt"
	"go/ast"
	"go/token"
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

func DetectParametersAndQuery(fn *ast.FuncDecl) ([]Parameter, error) {
	var parameters []Parameter

	// Ensure the function has a body
	if fn.Body == nil {
		return nil, fmt.Errorf("function body is nil")
	}

	// Inspect the function body to find parameter and query accesses
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		// Check for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if the function being called is a Gin context method
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Detect path parameters (e.g., c.Param("id"))
		if selExpr.Sel.Name == "Param" && len(callExpr.Args) == 1 {
			arg, ok := callExpr.Args[0].(*ast.BasicLit)
			if ok && arg.Kind == token.STRING {
				paramName := strings.Trim(arg.Value, "\"")
				parameters = append(parameters, Parameter{
					Name:        paramName,
					In:          "path",
					Required:    true,
					Schema:      "string",
					Description: fmt.Sprintf("Path parameter '%s'", paramName),
				})
			}
		}

		// Detect query parameters (e.g., c.Query("name"))
		if selExpr.Sel.Name == "Query" && len(callExpr.Args) == 1 {
			arg, ok := callExpr.Args[0].(*ast.BasicLit)
			if ok && arg.Kind == token.STRING {
				queryName := strings.Trim(arg.Value, "\"")
				parameters = append(parameters, Parameter{
					Name:        queryName,
					In:          "query",
					Required:    false,
					Schema:      "string",
					Description: fmt.Sprintf("Query parameter '%s'", queryName),
				})
			}
		}

		// Detect headers (e.g., c.GetHeader("X-Correlation-ID"))
		if selExpr.Sel.Name == "GetHeader" && len(callExpr.Args) == 1 {
			arg, ok := callExpr.Args[0].(*ast.BasicLit)
			if ok && arg.Kind == token.STRING {
				headerName := strings.Trim(arg.Value, "\"")
				parameters = append(parameters, Parameter{
					Name:        headerName,
					In:          "header",
					Required:    false,
					Schema:      "string",
					Description: fmt.Sprintf("Header '%s'", headerName),
				})
			}
		}

		return true
	})

	if len(parameters) > 0 {
		return parameters, nil
	}
	return nil, fmt.Errorf("no parameters or query strings found")
}

func DetectHeaders(fn *ast.FuncDecl) ([]Header, error) {
	var headers []Header

	// Ensure the function has a body
	if fn.Body == nil {
		return nil, fmt.Errorf("function body is nil")
	}

	// Inspect the function body to find header accesses
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		// Check for call expressions
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if the function being called is c.Header
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || selExpr.Sel.Name != "Header" {
			return true
		}

		// Extract the header name
		if len(callExpr.Args) == 2 {
			arg, ok := callExpr.Args[0].(*ast.BasicLit)
			if ok && arg.Kind == token.STRING {
				headerName := strings.Trim(arg.Value, "\"")
				headers = append(headers, Header{
					StatusCode:  "200",
					Name:        headerName,
					Type:        "string",
					Required:    true,
					Description: fmt.Sprintf("Header '%s'", headerName),
				})
			}
		}

		return true
	})

	if len(headers) > 0 {
		return headers, nil
	}
	return nil, fmt.Errorf("no headers found")
}
