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
- 🧪 Auto schema generation from Go structs
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
Step 1: Annotate your handlers
```go
// @Summary Get user by ID
// @Description Returns a user by ID
// @Tags user
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Header 200 X-RateLimit-Remaining string true "Remaining quota"
// @Router /user/{id} [get]
func GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	c.Header("X-RateLimit-Remaining", "29")
	c.JSON(200, UserResponse{ID: id, Name: "John Doe"})
	c.JSON(400, ErrorResponse{Message: err.Error()})
}
```

---

Step 2: Generate the spec via CLI
```bash
openapi3gen generate --source ./examples --output ./swagger/openapi.json
```

---

Step 3: Serve Swagger UI
```go
r := gin.Default()
swagger.RegisterSwaggerUI(r, "/swagger", "./swagger/openapi.json")
```
Access at: http://localhost:8080/swagger

---

## 🗂️ Annotation Cheatsheet
| Annotation          | Purpose                                 |
| ------------------- | --------------------------------------- |
| `@GlobalTitle`  	  | Title for the about the APIs/Service    |
| `@GlobalVersion`    | Version of the swagger doc              |
| `@GlobalDescription`| Detailed explanation about APIs/Service |
| `@Summary`          | One-line summary                        |
| `@Description`      | Detailed endpoint explanation           |
| `@Tags`             | Group endpoints                         |
| `@Param`            | Parameters in `path`, `query`, `header` |
| `@RequestBody`      | JSON body payload with struct           |
| `@Success`          | Success Response code and return object |
| `@Failure`          | Failure Response code and return object |
| `@Header`           | Adds response header details            |
| `@Deprecated`       | Flags the route as deprecated in spec   |

---

## 🛠 Developer Notes
Run CLI locally
```bash
go run main.go generate --source ./examples --output ./swagger/openapi.json
```

Build CLI binary
```bash
go build -o openapi3gen main.go
```

---

## 📌 Roadmap
- ✅  OpenAPI 3.0 support

- ✅  Schema generation from Go structs

- ✅  Swagger UI

- ✅  CLI for static spec generation

- ✅  Auto Detection of Path, Query & Header parameter

- ✅  Auto Detection of Header

- ⌛ Security schemes (@Security)

- ⌛ Support enums, examples

- ⌛ JSON/YAML toggles

- ⌛ Support other golang web frameworks like echo, chi etc.

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




