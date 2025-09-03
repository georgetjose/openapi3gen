// @GlobalTitle My Service API
// @GlobalVersion 1.0.0
// @GlobalDescription This is a sample API for demonstrating OpenAPI generation with Gin and annotations.
package main

import (
	"log"

	"github.com/georgetjose/openapi3gen/pkg/generator"
	"github.com/georgetjose/openapi3gen/pkg/parser"
	"github.com/georgetjose/openapi3gen/pkg/ui"
	"github.com/gin-gonic/gin"
)

// @Summary Hello greeting
// @Description This endpoint is a sample.
// @Tags hello
// @Success 200 "Success"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 "Unauthorized"
// @Router /hello [get]
func HelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, world!"})
}

// @Summary Legacy greeting
// @Description This endpoint is deprecated
// @Tags hello
// @Deprecated
// @Success 200 "Legacy greeting response"
// @Router /hello-legacy [get]
func LegacyHello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "This is deprecated"})
}

// @Summary Get user by ID
// @Description Returns user data based on ID
// @Tags user
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "Returns the user object with id and name"
// @Header 200 X-RateLimit-Remaining string true "Remaining quota"
// @Router /user/{id} [get]
func GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	c.Header("X-RateLimit-Remaining", "29")
	c.JSON(200, UserResponse{ID: id, Name: "George T Jose"})
}

// @Summary Search user by name
// @Description Returns user data based on query param
// @Tags user
// @Param name query string true "Name of the user to search"
// @Param X-Correlation-ID header string false "Tracking ID for the request"
// @Success 200 {object} UserResponse "Returns the user object"
// @Header 200 X-RateLimit-Remaining string true "Remaining quota"
// @Router /user/search [get]
func SearchUserHandler(c *gin.Context) {
	name := c.Query("name")
	correlationID := c.GetHeader("X-Correlation-ID")

	c.Header("X-RateLimit-Remaining", "28")
	c.JSON(200, UserResponse{
		ID:   correlationID,
		Name: name,
	})
}

// @Summary Create a user
// @Description Creates a new user
// @Tags user
// @RequestBody {object} CreateUserRequest true "User payload"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Security BearerAuth
// @Router /users [post]
func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(201, UserResponse{ID: "123", Name: req.Name})
}

// @Summary Create a user Auto Detect
// @Description Creates a new user
// @Tags user
// @Router /usersauto [post]
func CreateUserHandlerAutoDetect(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(201, UserResponse{ID: "123", Name: req.Name})
}

// @Summary Search user by name
// @Description Returns user data based on query param
// @Tags user
// @Router /user/searchauto [get]
func SearchUserHandlerAuto(c *gin.Context) {
	name := c.Query("name")
	correlationID := c.GetHeader("X-Correlation-ID")

	c.Header("X-RateLimit-Remaining", "28")
	c.JSON(200, UserResponse{
		ID:   correlationID,
		Name: name,
	})
}

type CreateUserRequest struct {
	Name    string  `json:"name" openapi:"desc=Full name of the user"`
	Email   string  `json:"email" openapi:"desc=User's email address"`
	Address Address `json:"address" openapi:"desc=User's address"`
}

type Address struct {
	State   string `json:"state" openapi:"desc=State"`
	ZipCode int    `json:"zip_code" openapi:"desc=ZIP code"`
}

type UserResponse struct {
	ID          string      `json:"id" openapi:"desc=Unique user ID"`
	Name        string      `json:"name" openapi:"desc=Full name of the user"`
	Description Description `json:"description" openapi:"desc=User description"`
}

type Description struct {
	Status  string `json:"status" openapi:"desc=User status"`
	Message string `json:"message" openapi:"desc=User's message"`
}

type ErrorResponse struct {
	Message string `json:"message" openapi:"desc=Error message"`
}

type HelloResponse struct {
	Message string `json:"message" openapi:"desc=Hello message"`
}

func main() {
	r := gin.Default()

	// Parse annotations
	routes, err := parser.ParseDirectory("./")
	if err != nil {
		log.Fatal(err)
	}

	registry := generator.NewModelRegistry()
	registry.Register("CreateUserRequest", CreateUserRequest{})
	registry.Register("UserResponse", UserResponse{})
	registry.Register("ErrorResponse", ErrorResponse{})
	registry.Register("Description", Description{})
	registry.Register("HelloResponse", HelloResponse{})

	globalMetaData := parser.ParseGlobalMetadata("main.go")
	// Generate OpenAPI spec
	openapi := generator.GenerateSpec(routes, registry, globalMetaData)

	// Serve OpenAPI dynamically
	ui.RegisterSwaggerUI(r, "")
	ui.RegisterSwaggerJSONHandler(r, openapi)

	// Routes available
	r.GET("/hello", HelloHandler)

	r.GET("/hello-legacy", LegacyHello)

	r.GET("/user/:id", GetUserByIDHandler)

	r.GET("/user/search", SearchUserHandler)

	r.POST("/users", CreateUserHandler)

	r.POST("/usersauto", CreateUserHandlerAutoDetect)

	r.GET("/user/searchauto", SearchUserHandlerAuto)

	r.Run(":8081")
}
