package generator

import (
	"fmt"
	"log"
	"strings"

	"github.com/georgetjose/openapi3gen/pkg/parser"
)

// GenerateSpec builds an OpenAPI struct from parsed RouteDoc list
func GenerateSpec(routes []parser.RouteDoc, registry *ModelRegistry, globalMetaData parser.GlobalMetadata) *OpenAPI {
	openapi := &OpenAPI{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       globalMetaData.GlobalTitle,
			Version:     globalMetaData.GlobalVersion,
			Description: globalMetaData.GlobalDescription,
		},
		Paths: make(map[string]*PathItem),
	}
	openapi.Components = &Components{
		Schemas: make(map[string]*Schema),
	}

	for _, route := range routes {
		pathItem, exists := openapi.Paths[route.Path]
		if !exists {
			pathItem = &PathItem{}
			openapi.Paths[route.Path] = pathItem
		}

		// 🔹 Reset parameters and requestBody per route
		var parameters []*ParameterObject
		var requestBody *RequestBodyObject

		// 🔹 Deduplication map
		seenParams := make(map[string]bool)

		for _, p := range route.Params {
			paramKey := p.In + ":" + p.Name
			if seenParams[paramKey] {
				continue // skip duplicate
			}
			seenParams[paramKey] = true

			if p.In == "path" && !strings.Contains(route.Path, "{"+p.Name+"}") {
				fmt.Printf("Warning: Path param '%s' not found in route path '%s'. Skipping.\n", p.Name, route.Path)
				continue
			}

			parameters = append(parameters, &ParameterObject{
				Name:        p.Name,
				In:          p.In,
				Required:    p.Required,
				Description: p.Description,
				Schema: &Schema{
					Type: p.Schema,
				},
			})
		}

		if route.RequestBody != nil {
			if m, ok := registry.Get(route.RequestBody.Model); ok {
				refSchema := addComponentSchema(route.RequestBody.Model, m, registry, openapi.Components)

				requestBody = &RequestBodyObject{
					Description: route.RequestBody.Description,
					Required:    route.RequestBody.Required,
					Content: map[string]MediaType{
						route.RequestBody.MediaType: {
							Schema: refSchema,
						},
					},
				}
			} else {
				log.Printf("Model not found in registry: %s\n", route.RequestBody.Model)
			}
		}

		responses := make(map[string]*ResponseWrapper)
		for statusCode, r := range route.Responses {
			if m, ok := registry.Get(r.Model); ok {
				refSchema := addComponentSchema(r.Model, m, registry, openapi.Components)

				// Collect response headers
				headers := make(map[string]*HeaderObject)
				for _, h := range route.Headers {
					if h.StatusCode == statusCode {
						headers[h.Name] = &HeaderObject{
							Description: h.Description,
							Schema: &Schema{
								Type: h.Type,
							},
						}
					}
				}

				desc := r.Description
				if desc == "" {
					desc = "Response"
				}

				// Add the response to the map without overwriting existing entries
				responses[statusCode] = &ResponseWrapper{
					Description: desc,
					Content: map[string]MediaType{
						r.MediaType: {
							Schema: refSchema,
						},
					},
					Headers: headers,
				}
			}
		}

		// Ensure at least 1 response
		if len(responses) == 0 {
			responses["200"] = &ResponseWrapper{
				Description: "OK",
			}
		}

		op := &Operation{
			Summary:     route.Summary,
			Description: route.Description,
			Tags:        route.Tags,
			Parameters:  parameters,
			RequestBody: requestBody,
			Responses:   responses,
		}

		if route.Deprecated {
			op.Deprecated = true
		}

		switch strings.ToLower(route.Method) {
		case "get":
			pathItem.Get = op
		case "post":
			pathItem.Post = op
		case "put":
			pathItem.Put = op
		case "delete":
			pathItem.Delete = op
		}
	}

	return openapi
}

func addComponentSchema(modelName string, model any, registry *ModelRegistry, components *Components) *Schema {
	// If already registered, return $ref
	if _, exists := components.Schemas[modelName]; exists {
		return &Schema{
			Ref: "#/components/schemas/" + modelName,
		}
	}

	schema := GenerateSchemaFromStruct(model)
	components.Schemas[modelName] = schema

	return &Schema{
		Ref: "#/components/schemas/" + modelName,
	}
}
