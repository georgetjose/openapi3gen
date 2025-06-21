package generator

type Schema struct {
	Type        string             `json:"type,omitempty" yaml:"type,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty" yaml:"items,omitempty"`
	Ref         string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
}

type Components struct {
	Schemas         map[string]*Schema               `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Ref             string                           `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeObject `json:"securitySchemes,omitempty"`
}

type OpenAPI struct {
	OpenAPI    string               `json:"openapi" yaml:"openapi"`
	Info       Info                 `json:"info" yaml:"info"`
	Paths      map[string]*PathItem `json:"paths" yaml:"paths"`
	Components *Components          `json:"components,omitempty" yaml:"components,omitempty"`
}

type Info struct {
	Title       string `json:"title" yaml:"title"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Post   *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Put    *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
}

type SecuritySchemeObject struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

type Operation struct {
	Summary     string                      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                      `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        []string                    `json:"tags,omitempty" yaml:"tags,omitempty"`
	Responses   map[string]*ResponseWrapper `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters  []*ParameterObject          `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBodyObject          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Security    []map[string][]string       `json:"security,omitempty"`
	Deprecated  bool                        `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

type ResponseWrapper struct {
	Description string                   `json:"description" yaml:"description"`
	Content     map[string]MediaType     `json:"content,omitempty" yaml:"content,omitempty"`
	Headers     map[string]*HeaderObject `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type MediaType struct {
	Schema *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type ParameterObject struct {
	Name        string  `json:"name" yaml:"name"`
	In          string  `json:"in" yaml:"in"`
	Required    bool    `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
}

type RequestBodyObject struct {
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool                 `json:"required,omitempty" yaml:"required,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

type HeaderObject struct {
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}
