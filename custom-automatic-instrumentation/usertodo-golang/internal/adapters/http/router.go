package http

import (
	"errors"
	"net/http"
	"time"

	"interview/internal/application/usecase"
	"interview/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createUserUseCase *usecase.CreateUserUseCase
	createTaskUseCase *usecase.CreateTaskUseCase
	logger            *logging.Logger
}

type createUserRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	TaskIDs []uint `json:"task_ids"`
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func DisableDebugMode() {
	gin.SetMode(gin.ReleaseMode)
}

func NewRouter(logger *logging.Logger, createUserUseCase *usecase.CreateUserUseCase, createTaskUseCase *usecase.CreateTaskUseCase) *gin.Engine {
	handler := &Handler{
		createUserUseCase: createUserUseCase,
		createTaskUseCase: createTaskUseCase,
		logger:            logger.WithModule("router"),
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger.WithModule("http")))
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/tasks", handler.createTask)
	router.POST("/users", handler.createUser)

	return router
}

func (h *Handler) createTask(ctx *gin.Context) {
	var request createTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid create task request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.createTaskUseCase.Execute(ctx.Request.Context(), usecase.CreateTaskInput{
		Title:       request.Title,
		Description: request.Description,
	})
	if err != nil {
		h.logger.Warn("create task failed title=%s error=%v", request.Title, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("task created id=%d title=%s done=%t", task.ID, task.Title, task.Done)

	ctx.JSON(http.StatusCreated, gin.H{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"done":        task.Done,
		"created_at":  task.CreatedAt,
	})
}

func (h *Handler) createUser(ctx *gin.Context) {
	var request createUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid create user request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.createUserUseCase.Execute(ctx.Request.Context(), usecase.CreateUserInput{
		Name:    request.Name,
		Email:   request.Email,
		TaskIDs: request.TaskIDs,
	})
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, usecase.ErrTaskNotFound) {
			statusCode = http.StatusNotFound
		}

		h.logger.Warn("create user failed email=%s task_ids=%v error=%v", request.Email, request.TaskIDs, err)
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	tasks := make([]gin.H, 0, len(user.Tasks))
	for _, task := range user.Tasks {
		tasks = append(tasks, gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"done":        task.Done,
		})
	}

	h.logger.Info("user created id=%d email=%s task_count=%d", user.ID, user.Email, len(user.Tasks))

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"tasks":      tasks,
		"created_at": user.CreatedAt,
	})
}

func requestLogger(logger *logging.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		logger.Info("%s %s status=%d latency_ms=%d client_ip=%s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Writer.Status(),
			time.Since(start).Milliseconds(),
			ctx.ClientIP(),
		)
	}
}
