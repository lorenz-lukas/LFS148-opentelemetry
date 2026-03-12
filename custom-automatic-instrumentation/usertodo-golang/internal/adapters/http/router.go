package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"interview/internal/application/usecase"
	"interview/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createUserUseCase       *usecase.CreateUserUseCase
	listTasksUseCase        *usecase.ListTasksUseCase
	attachTaskToUserUseCase *usecase.AttachTaskToUserUseCase
	logger                  *logging.Logger
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// swagger:response errorResponse
type errorResponse struct {
	// in: body
	Body struct {
		Error string `json:"error"`
	}
}

// swagger:response healthResponse
type healthResponse struct {
	// in: body
	Body struct {
		Status string `json:"status"`
	}
}

// swagger:response todoResponse
type todoResponse struct {
	// in: body
	Body todoResponseDoc
}

// swagger:response todosResponse
type todosResponse struct {
	// in: body
	Body []todoResponseDoc
}

// swagger:response userResponse
type userResponse struct {
	// in: body
	Body userResponseDoc
}

// swagger:model CreateUserRequest
type createUserRequestDoc struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// swagger:model TodoResponse
type todoResponseDoc struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

// swagger:model UserResponse
type userResponseDoc struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Tasks     []todoResponseDoc `json:"tasks"`
	CreatedAt time.Time         `json:"created_at"`
}

// swagger:parameters createUser
type createUserParamsWrapper struct {
	// in: body
	// required: true
	Body createUserRequestDoc
}

// swagger:parameters attachTaskToUser
type attachTaskToUserParamsWrapper struct {
	// in: path
	// required: true
	UserID uint `json:"userId"`
	// in: path
	// required: true
	TaskID uint `json:"taskId"`
}

func DisableDebugMode() {
	gin.SetMode(gin.ReleaseMode)
}

func NewRouter(
	logger *logging.Logger,
	createUserUseCase *usecase.CreateUserUseCase,
	listTasksUseCase *usecase.ListTasksUseCase,
	attachTaskToUserUseCase *usecase.AttachTaskToUserUseCase,
) *gin.Engine {
	handler := &Handler{
		createUserUseCase:       createUserUseCase,
		listTasksUseCase:        listTasksUseCase,
		attachTaskToUserUseCase: attachTaskToUserUseCase,
		logger:                  logger.WithModule("router"),
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger.WithModule("http")))

	// Serve the checked-in Swagger document and a minimal UI page.
	router.GET("/swagger", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})
	router.GET("/swagger/doc.yaml", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "application/yaml; charset=utf-8", swaggerYAML)
	})

	// swagger:route GET /health health health
	//
	// Service health check.
	//
	// Responses:
	//   200: healthResponse
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/todos", handler.listTodos)
	router.POST("/users", handler.createUser)
	router.POST("/users/:userID/tasks/:taskID", handler.attachTaskToUser)

	return router
}

// swagger:route GET /todos todos listTodos
//
// List todos.
//
// Responses:
//   200: todosResponse
//   500: errorResponse
func (h *Handler) listTodos(ctx *gin.Context) {
	tasks, err := h.listTasksUseCase.Execute(ctx.Request.Context())
	if err != nil {
		h.logger.Warn("list todos failed error=%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]gin.H, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"done":        task.Done,
			"created_at":  task.CreatedAt,
		})
	}

	h.logger.Info("todos listed count=%d", len(tasks))
	ctx.JSON(http.StatusOK, response)
}

// swagger:route POST /users users createUser
//
// Create a user.
//
// Responses:
//   201: userResponse
//   400: errorResponse
func (h *Handler) createUser(ctx *gin.Context) {
	var request createUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid create user request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.createUserUseCase.Execute(ctx.Request.Context(), usecase.CreateUserInput{
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		h.logger.Warn("create user failed email=%s error=%v", request.Email, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("user created id=%d email=%s", user.ID, user.Email)

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"tasks":      []gin.H{},
		"created_at": user.CreatedAt,
	})
}

// swagger:route POST /users/{userId}/tasks/{taskId} users attachTaskToUser
//
// Link an existing todo to an existing user.
//
// Responses:
//   200: userResponse
//   400: errorResponse
//   404: errorResponse
func (h *Handler) attachTaskToUser(ctx *gin.Context) {
	userID, err := parseUintParam(ctx.Param("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	taskID, err := parseUintParam(ctx.Param("taskID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	user, err := h.attachTaskToUserUseCase.Execute(ctx.Request.Context(), usecase.AttachTaskToUserInput{
		UserID: userID,
		TaskID: taskID,
	})
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, usecase.ErrTaskNotFound) || errors.Is(err, usecase.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}

		h.logger.Warn("attach task to user failed user_id=%d task_id=%d error=%v", userID, taskID, err)
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
			"created_at":  task.CreatedAt,
		})
	}

	h.logger.Info("task attached user_id=%d task_id=%d", userID, taskID)
	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"tasks":      tasks,
		"created_at": user.CreatedAt,
	})
}

func parseUintParam(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(parsed), nil
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
