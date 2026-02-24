package main

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()
	userCollection := mongoose("users")
	todosCollection := mongoose("todos")

	r.POST("/login", func(ctx *gin.Context) {
		Login(ctx, userCollection)
	})

	r.POST("/register", func(ctx *gin.Context) {
		Register(ctx, userCollection)
	})

	protected := r.Group("/")
	protected.Use(AuthMiddleware())

	protected.GET("/todos", func(c *gin.Context) {
		GetTodos(c, todosCollection)
	})

	protected.GET("/todos/:id", func(c *gin.Context) {
		GetTodo(c, todosCollection)
	})

	protected.POST("/todos", func(c *gin.Context) {
		CreateTodo(c, todosCollection)
	})

	protected.DELETE("/todos/:id", func(c *gin.Context) {
		DeleteTodo(c, todosCollection)
	})

	return r
}