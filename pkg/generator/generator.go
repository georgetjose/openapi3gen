package generator

import (
	"fmt"
	"log"
	"reflect"
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
		Schemas:         make(map[string]*Schema),
		SecuritySchemes: make(map[string]*SecuritySchemeObject),
	}

	usedSecurity := make(map[string]bool)
	for _, route := range routes {
		for _, sec := range route.SecuritySchemes {
			if !usedSecurity[sec] {
				openapi.Components.SecuritySchemes[sec] = &SecuritySchemeObject{
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				}
				usedSecurity[sec] = true
			}
		}
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
				refSchema := addComponentSchema(route.RequestBody.Model, m, openapi.Components)

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

			response := &ResponseWrapper{
				Description: desc,
				Headers:     headers,
			}

			// Only add content if there's a model
			if r.Model != "" {
				if m, ok := registry.Get(r.Model); ok {
					refSchema := addComponentSchema(r.Model, m, openapi.Components)
					response.Content = map[string]MediaType{
						r.MediaType: {
							Schema: refSchema,
						},
					}
				} else {
					log.Printf("Model not found in registry: %s\n", r.Model)
					continue // Skip this response if model not found
				}
			}

			responses[statusCode] = response
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

		var security []map[string][]string
		for _, secName := range route.SecuritySchemes {
			security = append(security, map[string][]string{
				secName: {},
			})
		}
		op.Security = security

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

func addComponentSchema(modelName string, model any, components *Components) *Schema {
	// If already registered, return $ref
	if _, exists := components.Schemas[modelName]; exists {
		return &Schema{
			Ref: "#/components/schemas/" + modelName,
		}
	}

	schema := GenerateSchemaFromStruct(model)
	components.Schemas[modelName] = schema

	// Recursively register nested struct schemas
	registerNestedSchemas(model, components)

	return &Schema{
		Ref: "#/components/schemas/" + modelName,
	}
}

// registerNestedSchemas recursively registers schemas for nested structs
func registerNestedSchemas(model any, components *Components) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// If it's a custom struct, register it
		if fieldType.Kind() == reflect.Struct && isCustomStruct(fieldType) {
			schemaName := fieldType.Name()

			// If not already registered, register it
			if _, exists := components.Schemas[schemaName]; !exists {
				// Create an instance of the struct to generate schema
				structValue := reflect.New(fieldType).Interface()
				schema := GenerateSchemaFromStruct(structValue)
				components.Schemas[schemaName] = schema

				// Recursively register nested structs
				registerNestedSchemas(structValue, components)
			}
		}
	}
}
