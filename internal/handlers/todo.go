package handlers

import (
	"log"
	"net/http"
	"strconv"
	"todo-list-api/internal/models"
	"todo-list-api/internal/services"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new todo
// @Description Create a new todo item (requires authentication)
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.Todo true "Todo item"
// @Success 201 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /todos [post]
func CreateTodo(todoService services.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo models.Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.MustGet("user_id").(int)
		go func() {
			log.Printf("User %d is creating a todo: %s", userID, todo.Title)
		}() // Goroutine para log assíncrono
		if err := todoService.CreateTodo(&todo, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, todo)
	}
}

// @Summary Get all todos
// @Description Get all todo items (public endpoint)
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Router /todos [get]
func GetAllTodos(todoService services.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := todoService.GetAllTodos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}

// @Summary Update a todo
// @Description Update a todo item (requires authentication, only owner)
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Updated todo"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /todos/{id} [put]
func UpdateTodo(todoService services.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var todo models.Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		todo.ID = id
		userID := c.MustGet("user_id").(int)
		if err := todoService.UpdateTodo(&todo, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

// @Summary Delete a todo
// @Description Delete a todo item (requires authentication, only owner)
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /todos/{id} [delete]
func DeleteTodo(todoService services.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		userID := c.MustGet("user_id").(int)
		if err := todoService.DeleteTodo(id, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}