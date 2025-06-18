package main

import (
	"log"

	"github.com/georgetjose/openapi3gen/pkg/generator"
	"github.com/georgetjose/openapi3gen/pkg/parser"
	"github.com/georgetjose/openapi3gen/pkg/ui"
	"github.com/gin-gonic/gin"
)

// @Summary Greet user
// @Description Returns a friendly greeting message
// @Tags hello
// @Param name query string true "Name of the user"
// @Success 200 {object} map[string]string
// @Router /hello [get]
func HelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, world!"})
}

// @Summary Legacy greeting
// @Description This endpoint is deprecated
// @Tags legacy
// @Deprecated
// @Success 200 {object} map[string]string
// @Router /hello-legacy [get]
func LegacyHello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "This is deprecated"})
}

// @Summary Get user by ID
// @Description Returns user data based on ID
// @Tags user
// @Param id path string true "User ID"
// @Param X-Correlation-ID header string true "Tracking ID for the request"
// @Success 200 {object} UserResponse "Returns the user object with id and name"
// @Header 200 X-RateLimit-Remaining string true "Remaining quota"
// @Router /user/{id} [get]
func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	c.Header("X-RateLimit-Remaining", "29")
	c.JSON(200, UserResponse{ID: id, Name: "George T Jose"})
}

// @Summary Create a user
// @Description Creates a new user
// @Tags user
// @RequestBody {object} CreateUserRequest true "User payload"
// @Success 201 {object} UserResponse
// @Router /users [post]
func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, UserResponse{ID: "123", Name: req.Name})
}

type CreateUserRequest struct {
	Name  string `json:"name" openapi:"desc=Full name of the user"`
	Email string `json:"email" openapi:"desc=User's email address"`
}

type UserResponse struct {
	ID   string `json:"id" openapi:"desc=Unique user ID"`
	Name string `json:"name" openapi:"desc=Full name of the user"`
}

func main() {
	r := gin.Default()

	// Parse annotations
	routes, err := parser.ParseDirectory("./examples/gin-basic")
	if err != nil {
		log.Fatal(err)
	}

	registry := generator.NewModelRegistry()
	registry.Register("CreateUserRequest", CreateUserRequest{})
	registry.Register("UserResponse", UserResponse{})

	// Generate OpenAPI spec
	openapi := generator.GenerateSpec(routes, registry)

	// Serve OpenAPI dynamically
	ui.RegisterSwaggerUI(r, "")
	ui.RegisterSwaggerJSONHandler(r, openapi)

	// Routes
	r.GET("/hello", HelloHandler)

	r.GET("/hello-legacy", LegacyHello)

	r.GET("/user/:id", GetUserHandler)

	r.POST("/users", CreateUserHandler)

	r.Run(":8080")
}
