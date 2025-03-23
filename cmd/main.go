package main

import (
	"log"
	"todo-list-api/internal/handlers"
	"todo-list-api/internal/middleware"
	"todo-list-api/internal/repositories"
	"todo-list-api/internal/services"
	"todo-list-api/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "todo-list-api/docs" 
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	// Repositórios
	userRepo := repositories.NewUserRepository(database.DB)
	todoRepo := repositories.NewTodoRepository(database.DB)

	// Serviços
	authService := services.NewAuthService(userRepo)
	todoService := services.NewTodoService(todoRepo)

	r := gin.Default()

	// Rotas de autenticação
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", handlers.Register(authService))
		authGroup.POST("/login", handlers.Login(authService))
		authGroup.POST("/refresh", handlers.Refresh(authService))
	}

	// Rotas de to-dos
	todoGroup := r.Group("/todos")
	{
		todoGroup.GET("", handlers.GetAllTodos(todoService))
		todoGroup.Use(middleware.AuthMiddleware())
		todoGroup.POST("", handlers.CreateTodo(todoService))
		todoGroup.PUT("/:id", handlers.UpdateTodo(todoService))
		todoGroup.DELETE("/:id", handlers.DeleteTodo(todoService))
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Inicia o servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}