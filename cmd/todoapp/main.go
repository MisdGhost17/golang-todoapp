package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/MisdGhost17/golang-todoapp/internal/core/config"
	core_logger "github.com/MisdGhost17/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/MisdGhost17/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/MisdGhost17/golang-todoapp/internal/core/transport/http/server"
	auth_passwordhash "github.com/MisdGhost17/golang-todoapp/internal/features/auth/password_hash"
	auth_postgres_repository "github.com/MisdGhost17/golang-todoapp/internal/features/auth/repository/postgres"
	auth_service "github.com/MisdGhost17/golang-todoapp/internal/features/auth/service"
	auth_token "github.com/MisdGhost17/golang-todoapp/internal/features/auth/token"
	auth_transport_http "github.com/MisdGhost17/golang-todoapp/internal/features/auth/transport/http"
	statistics_postgres_repository "github.com/MisdGhost17/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/MisdGhost17/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/MisdGhost17/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/MisdGhost17/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/MisdGhost17/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/MisdGhost17/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/MisdGhost17/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/MisdGhost17/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/MisdGhost17/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/MisdGhost17/golang-todoapp/docs"
)

// @title 		 Golang ToDo API
// @version 	 1.0
// @description  ToDo Application REST API
// @host 		 127.0.0.1:5050
// @BasePath 	 /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите Bearer <access_token>. Токен можно получить через POST /auth/login.
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)

	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	taskTransportHttp := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHttp := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	tokenManager, err := auth_token.NewJWTManager(
		[]byte(os.Getenv("JWT_SECRET")),
		"todoapp",
		15*time.Minute,
	)
	if err != nil {
		logger.Fatal("failed to initialize JWT manager", zap.Error(err))
	}

	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	passwordHasher := auth_passwordhash.NewBcryptHasher()
	refreshTokenManager := auth_token.NewRefreshTokenManager()

	authService := auth_service.NewAuthService(
		authRepository,
		passwordHasher,
		tokenManager,
		refreshTokenManager,
	)

	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHttp.Routes()...)
	apiVersionRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHttp.Routes()...)
	requireAuth := core_http_middleware.RequireAuth(tokenManager)

	taskRoutes := taskTransportHttp.Routes()
	for i := range taskRoutes {
		taskRoutes[i].Middleware = append(
			[]core_http_middleware.Middleware{requireAuth},
			taskRoutes[i].Middleware...,
		)
	}

	apiVersionRouterV1.RegisterRoutes(taskRoutes...)

	// apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
	// 	core_http_server.ApiVersion2,
	// 	core_http_middleware.Dummy("api v2 middleware"),
	// )
	// apiVersionRouterV2.RegisterRoutes(usersTransportHttp.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouterV1)
	// httpServer.RegisterAPIRouters(apiVersionRouterV2)
	httpServer.RegisterSwagger()
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
