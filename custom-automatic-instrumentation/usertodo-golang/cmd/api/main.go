package main

import (
	"os"

	httpadapter "interview/internal/adapters/http"
	postgresadapter "interview/internal/adapters/postgres"
	"interview/internal/application/usecase"
	"interview/internal/infrastructure/config"
	"interview/internal/infrastructure/database"
	"interview/internal/infrastructure/logging"
)

func main() {
	logger := logging.New("root", "main")
	httpadapter.DisableDebugMode()

	cfg := config.Load()
	baseURL := "http://localhost:" + cfg.Port
	logger.Info("starting api on %s", baseURL)

	db, err := database.NewPostgres(cfg.DatabaseDSN)
	if err != nil {
		logger.Error("connect database: %v", err)
		os.Exit(1)
	}
	logger.Info("database connected")

	if err := postgresadapter.AutoMigrate(db); err != nil {
		logger.Error("migrate database: %v", err)
		os.Exit(1)
	}
	logger.Info("database migration completed")

	userRepository := postgresadapter.NewUserRepository(db)
	taskRepository := postgresadapter.NewTaskRepository(db)

	createUserUseCase := usecase.NewCreateUserUseCase(userRepository, taskRepository)
	createTaskUseCase := usecase.NewCreateTaskUseCase(taskRepository)

	router := httpadapter.NewRouter(logger, createUserUseCase, createTaskUseCase)

	logger.Info("http server listening on %s", baseURL)
	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Error("run http server: %v", err)
		os.Exit(1)
	}
}
