# 🧬 openapi3gen

> A Code-First OpenAPI 3.0 Specification Generator for Go (Gin) with Swagger UI support

`openapi3gen` is a lightweight, code-first tool that parses annotations in your Gin HTTP handlers and generates a complete OpenAPI 3.0 spec — ready to be served live or viewed in Swagger UI.

---

## ✨ Features

- ✅ Supports **OpenAPI 3.0** compliant schema
- 📌 Generates `openapi.json` from handler annotations
- 🔍 Path, query, and header param support via `@Param`
- 📦 Request body model via `@RequestBody`
- 🧾 Response models with `$ref`, `@Success` and `@Failure`
- 📤 Response headers with `@Header`
- 🔐 Multiple security schemes: `BearerAuth`, `ApiKeyAuth` with custom headers
- 🤖 **Auto-detection** of parameters, request bodies, responses, and headers
- 🧪 Auto schema generation from Go structs with `openapi` tags
- 🏷 Tag-based grouping, descriptions, and `@Deprecated`
- ⚙️ CLI support: `openapi3gen generate`
- 🌐 Swagger UI integration via embedded static assets

---

## 📦 Installation

```bash
go get github.com/georgetjose/openapi3gen

```

---

## 🚀 Getting Started

### Step 1: Annotate your handlers
```go
// @Summary Get user by ID
// @Description Returns a user by ID
// @Tags user
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "Returns the user object with id and name"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Header 200 X-RateLimit-Remaining string true "Remaining quota"
// @Security ApiKeyAuth:X-User-Token
// @Router /user/{id} [get]
func GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	c.Header("X-RateLimit-Remaining", "29")
	c.JSON(200, UserResponse{ID: id, Name: "John Doe"})
	if err != nil {
		c.JSON(400, ErrorResponse{Message: err.Error()})
	}
}
```

### Step 2: Add global metadata (optional)
```go
// @GlobalTitle My Service API
// @GlobalVersion 1.0.0
// @GlobalDescription This is a sample API for demonstrating OpenAPI generation.
package main
```

---

### Step 3: Register your models
```go
registry := generator.NewModelRegistry()
registry.Register("UserResponse", UserResponse{})
registry.Register("ErrorResponse", ErrorResponse{})
```

---

### Step 4: Generate the spec via CLI
```bash
openapi3gen generate --dir ./examples --output ./swagger/openapi.json
```

---

### Step 5: Serve Swagger UI (optional)
```go
r := gin.Default()

// Parse and generate OpenAPI spec
routes, _ := parser.ParseDirectory("./")
registry := generator.NewModelRegistry()
registry.Register("UserResponse", UserResponse{})
globalMetaData := parser.ParseGlobalMetadata("main.go")
openapi := generator.GenerateSpec(routes, registry, globalMetaData)

// Serve Swagger UI and JSON
ui.RegisterSwaggerUI(r, "")
ui.RegisterSwaggerJSONHandler(r, openapi)
```
Access at: http://localhost:8080/swagger

---

## 🗂️ Annotation Cheatsheet

| Annotation             | Purpose                                        | Example |
| ---------------------- | ---------------------------------------------- | ------- |
| `@GlobalTitle`  	     | Title for the APIs/Service                     | `@GlobalTitle My Service API` |
| `@GlobalVersion`       | Version of the swagger doc                     | `@GlobalVersion 1.0.0` |
| `@GlobalDescription`   | Detailed explanation about APIs/Service        | `@GlobalDescription This is a sample API` |
| `@Summary`             | One-line summary                               | `@Summary Get user by ID` |
| `@Description`         | Detailed endpoint explanation                  | `@Description Returns user data based on ID` |
| `@Tags`                | Group endpoints                                | `@Tags user,admin` |
| `@Param`               | Parameters in `path`, `query`, `header`        | `@Param id path string true "User ID"` |
| `@RequestBody`         | JSON body payload with struct                  | `@RequestBody {object} UserRequest true "User data"` |
| `@Success`             | Success Response code and return object        | `@Success 200 {object} UserResponse "Success"` |
| `@Failure`             | Failure Response code and return object        | `@Failure 400 {object} ErrorResponse "Bad Request"` |
| `@Header`              | Adds response header details                   | `@Header 200 X-RateLimit string true "Rate limit"` |
| `@Security`            | Adds authorization to endpoints                | `@Security BearerAuth` or `@Security ApiKeyAuth:X-Token` |
| `@Deprecated`          | Flags the route as deprecated in spec          | `@Deprecated` |

## 🔐 Security Schemes

### Bearer Authentication
```go
// @Security BearerAuth
// @Router /protected [get]
func ProtectedHandler(c *gin.Context) { ... }
```

### API Key Authentication
```go
// Default header (X-API-Key)
// @Security ApiKeyAuth
// @Router /api [get]

// Custom header with colon notation
// @Security ApiKeyAuth:X-User-Token
// @Router /user-auth [get]

// Custom header with bracket notation
// @Security ApiKeyAuth[x-territory-Key]
// @Router /territory [get]
```

## 🤖 Auto-Detection Features

openapi3gen can automatically detect many elements from your code:

- **Parameters**: Detects `c.Param()`, `c.Query()`, `c.GetHeader()` calls
- **Request Bodies**: Detects `c.ShouldBindJSON()` usage
- **Response Bodies**: Detects `c.JSON()` calls with status codes
- **Response Headers**: Detects `c.Header()` calls

### Auto-Detection Example
```go
// Minimal annotations - most things are auto-detected
// @Summary Create user auto-detect
// @Description Creates a new user with auto-detection
// @Tags user
// @Router /users/auto [post]
func CreateUserAutoHandler(c *gin.Context) {
    var req CreateUserRequest  // Auto-detected as request body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Message: err.Error()})  // Auto-detected as 400 response
        return
    }
    c.JSON(201, UserResponse{ID: "123", Name: req.Name})  // Auto-detected as 201 response
}
```

## 🏗️ Struct Schema Generation

Use `openapi` struct tags for enhanced schema documentation:

```go
type CreateUserRequest struct {
    Name    string  `json:"name" openapi:"desc=Full name of the user"`
    Email   string  `json:"email" openapi:"desc=User's email address"`
    Address Address `json:"address" openapi:"desc=User's address"`
}

type Address struct {
    State   string `json:"state" openapi:"desc=State"`
    ZipCode int    `json:"zip_code" openapi:"desc=ZIP code"`
}
```

---

## 🛠 Developer Notes

### Run CLI locally
```bash
go run main.go generate --dir ./examples --output ./swagger/openapi.json
```

### Build CLI binary
```bash
go build -o openapi3gen main.go
```

### Model Registration
Remember to register all models referenced in annotations:
```go
registry := generator.NewModelRegistry()
registry.Register("CreateUserRequest", CreateUserRequest{})
registry.Register("UserResponse", UserResponse{})  
registry.Register("ErrorResponse", ErrorResponse{})
```

---

## 📌 Roadmap

- ✅  OpenAPI 3.0 support
- ✅  Schema generation from Go structs  
- ✅  Swagger UI
- ✅  CLI for static spec generation
- ✅  Auto Detection of Path, Query & Header parameters
- ✅  Auto Detection of Request Body
- ✅  Auto Detection of Response Headers
- ✅  Auto Detection of Response Body
- ✅  Support for multiple Security schemes (BearerAuth, ApiKeyAuth)
- ✅  Custom security headers with flexible notation
- ✅  Enhanced struct schema generation with `openapi` tags
- ✅  Nested struct support with automatic `$ref` generation
- ⌛ Support enums, examples
- ⌛ JSON/YAML output toggles
- ⌛ Support other golang web frameworks like echo, chi etc.
- ⌛ OpenAPI 3.1 support

---

## 🤝 Contributing
Contributions welcome!

🌟 Star the repo

🐛 File issues and suggestions

🧪 Add tests for new functionality

📥 Open PRs for features or fixes

---

## 📬 Contact
For questions, feedback, or ideas:

🤖 GitHub: @georgetjose

✉️ Email: georgeb4pc@gmail.com

---

## 📄 License
openapi3gen is released under the Apache 2.0 license




