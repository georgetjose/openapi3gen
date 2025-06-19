package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type GlobalMetadata struct {
	GlobalTitle       string
	GlobalVersion     string
	GlobalDescription string
}

type Parameter struct {
	Name        string
	In          string // path, query, header, cookie
	Required    bool
	Schema      string // string, integer, etc.
	Description string
}

type RequestBody struct {
	Model       string
	Required    bool
	MediaType   string // Default: application/json
	Description string
}

type Response struct {
	Model       string
	MediaType   string
	StatusCode  string
	Description string
}

type Header struct {
	StatusCode  string
	Name        string
	Type        string
	Required    bool
	Description string
}

type RouteDoc struct {
	Summary     string
	Description string
	Method      string
	Path        string
	Tags        []string
	Params      []Parameter
	RequestBody *RequestBody
	Responses   map[string]Response
	Headers     []Header
	Deprecated  bool
}

func ParseGlobalMetadata(filePath string) GlobalMetadata {
	src, _ := os.ReadFile(filePath)
	lines := strings.Split(string(src), "\n")

	metadata := GlobalMetadata{}
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "// @GlobalTitle ") {
			metadata.GlobalTitle = strings.TrimSpace(strings.TrimPrefix(line, "// @GlobalTitle "))
		}
		if strings.HasPrefix(line, "// @GlobalVersion ") {
			metadata.GlobalVersion = strings.TrimSpace(strings.TrimPrefix(line, "// @GlobalVersion "))
		}
		if strings.HasPrefix(line, "// @GlobalDescription ") {
			metadata.GlobalDescription = strings.TrimSpace(strings.TrimPrefix(line, "// @GlobalDescription "))
		}
		if strings.HasPrefix(line, "package ") {
			break // Stop after reaching package line
		}
	}
	return metadata
}

// ParseDirectory parses all .go files in a folder and extracts annotations
func ParseDirectory(dir string) ([]RouteDoc, error) {
	var routes []RouteDoc

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, f := range node.Decls {
			fn, ok := f.(*ast.FuncDecl)
			if !ok || fn.Doc == nil {
				continue
			}

			doc := RouteDoc{
				Responses: make(map[string]Response),
			}

			for _, comment := range fn.Doc.List {
				text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))

				switch {
				case strings.HasPrefix(text, "@Summary "):
					doc.Summary = strings.TrimPrefix(text, "@Summary ")
				case strings.HasPrefix(text, "@Description "):
					doc.Description = strings.TrimPrefix(text, "@Description ")
				case strings.HasPrefix(text, "@Tags "):
					doc.Tags = strings.Split(strings.TrimPrefix(text, "@Tags "), ",")
				case strings.HasPrefix(text, "@Success "):
					// Format: @Success 200 {object} ModelName "Description"
					parts := strings.Fields(text[len("@Success "):])
					if len(parts) >= 3 {
						resp := Response{
							StatusCode: parts[0],
							MediaType:  "application/json",
							Model:      parts[2],
						}
						if len(parts) > 3 {
							resp.Description = strings.Join(parts[3:], " ")
							resp.Description = strings.Trim(resp.Description, `"`) // remove quotes
						}
						doc.Responses[parts[0]] = resp
					}
				case strings.HasPrefix(text, "@Router "):
					parts := strings.Fields(strings.TrimPrefix(text, "@Router "))
					if len(parts) == 2 {
						doc.Path = normalizePath(parts[0])
						doc.Method = strings.Trim(parts[1], "[]")
					}
				case strings.HasPrefix(text, "@Param "):
					// Format: @Param name in type required "description"
					// Example: @Param X-Correlation-ID header string true "Tracking ID"
					parts := strings.Fields(text[len("@Param "):])
					if len(parts) >= 4 {
						param := Parameter{
							Name:     parts[0],
							In:       parts[1],
							Schema:   parts[2],
							Required: parts[3] == "true",
						}
						if len(parts) > 4 {
							param.Description = strings.Join(parts[4:], " ")
							param.Description = strings.Trim(param.Description, `"`)
						}
						doc.Params = append(doc.Params, param)
					}
				case strings.HasPrefix(text, "@RequestBody "):
					// Format: @RequestBody {object} ModelName true "Description"
					parts := strings.Fields(text[len("@RequestBody "):])
					if len(parts) >= 4 {
						doc.RequestBody = &RequestBody{
							Model:       parts[1],                     // e.g., MyStruct
							Required:    parts[2] == "true",           // true or false
							Description: strings.Join(parts[3:], " "), // "User payload"
							MediaType:   "application/json",           // default for now
						}
					}
				case strings.HasPrefix(text, "@Header "):
					// Format: @Header 200 X-Header string true "Description"
					parts := strings.Fields(text[len("@Header "):])
					if len(parts) >= 5 {
						doc.Headers = append(doc.Headers, Header{
							StatusCode:  parts[0],
							Name:        parts[1],
							Type:        parts[2],
							Required:    parts[3] == "true",
							Description: strings.Join(parts[4:], " "),
						})
					}
				case strings.EqualFold(text, "@Deprecated"):
					doc.Deprecated = true
				}
			}

			if doc.Path != "" && doc.Method != "" {
				routes = append(routes, doc)
			}
		}

		return nil
	})

	return routes, err
}

func normalizePath(path string) string {
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			segments[i] = "{" + strings.TrimPrefix(segment, ":") + "}"
		}
	}
	return strings.Join(segments, "/")
}
