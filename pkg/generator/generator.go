package generator

import (
	"strings"

	"github.com/georgetjose/openapi3gen/pkg/parser"
)

// GenerateSpec builds an OpenAPI struct from parsed RouteDoc list
func GenerateSpec(routes []parser.RouteDoc, registry *ModelRegistry) *OpenAPI {
	openapi := &OpenAPI{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:   "Generated API",
			Version: "1.0.0",
		},
		Paths: make(map[string]*PathItem),
	}
	openapi.Components = &Components{
		Schemas: make(map[string]*Schema),
	}

	var parameters []*ParameterObject
	var requestBody *RequestBodyObject

	for _, route := range routes {
		pathItem, exists := openapi.Paths[route.Path]
		if !exists {
			pathItem = &PathItem{}
			openapi.Paths[route.Path] = pathItem
		}

		for _, p := range route.Params {
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
			}
		}

		responses := make(map[string]*ResponseWrapper)
		for statusCode, r := range route.Responses {
			if m, ok := registry.Get(r.Model); ok {
				refSchema := addComponentSchema(r.Model, m, registry, openapi.Components)

				// Collect all headers for this response
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
					desc = "Success"
				}
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

		op := &Operation{
			Summary:     route.Summary,
			Description: route.Description,
			Tags:        route.Tags,
			Parameters:  parameters,
			RequestBody: requestBody,
			Responses:   responses,
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
	// If already registered, just return the $ref
	if _, exists := components.Schemas[modelName]; exists {
		return &Schema{
			Ref: "#/components/schemas/" + modelName,
		}
	}

	// Generate and register schema
	schema := GenerateSchemaFromStruct(model)
	components.Schemas[modelName] = schema

	return &Schema{
		Ref: "#/components/schemas/" + modelName,
	}
}
