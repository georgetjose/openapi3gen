package generator

import (
	"reflect"
	"strings"
)

func GenerateSchemaFromStruct(model any) *Schema {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	schema := &Schema{
		Type:       "object",
		Properties: map[string]*Schema{},
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonName := field.Tag.Get("json")
		if jsonName == "" || jsonName == "-" {
			continue
		}

		// Remove ,omitempty etc.
		jsonName = parseJSONName(jsonName)

		desc := extractDescription(field.Tag.Get("openapi"))
		prop := &Schema{
			Type:        mapGoTypeToOpenAPIType(field.Type.Kind()),
			Description: desc,
		}

		schema.Properties[jsonName] = prop
	}

	return schema
}

func parseJSONName(tag string) string {
	if idx := len(tag); idx > 0 {
		if i := indexComma(tag); i >= 0 {
			return tag[:i]
		}
	}
	return tag
}

func indexComma(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}

func mapGoTypeToOpenAPIType(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	default:
		return "string"
	}
}

func extractDescription(tag string) string {
	// Example tag: openapi:"desc=User's name"
	parts := strings.Split(tag, "desc=")
	if len(parts) == 2 {
		return strings.Trim(parts[1], `"`)
	}
	return ""
}
