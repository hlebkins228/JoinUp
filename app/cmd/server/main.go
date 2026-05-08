package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"JoinUp/internal/controller"
	"JoinUp/internal/controller/handlers"
	"JoinUp/internal/repository"
	"JoinUp/internal/service"
	"JoinUp/internal/settings"
	"JoinUp/internal/utils/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// @title JoinUp API
// @version 1.0
// @description JoinUp backend API documentation.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

/*
TODO:
Добавить пакет dto в controller и на каждый запрос/ответ (почти на каждый) своя dto
Передалать некоторую логику с учетом новых dto
Сделать PUT запрос с path параметром в виде id
*/

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer ctxCancel()

	cfg, err := settings.ReadConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	var loggerCfg zap.Config
	switch cfg.Logging.Type {
	case settings.LoggerTypeProd:
		loggerCfg = zap.NewProductionConfig()
	case settings.LoggerTypeDev:
		loggerCfg = zap.NewDevelopmentConfig()
	default:
		panic("invalid logger type")
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		logger.Fatal("SERVER_PORT is not set")
	}
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		logger.Fatal("SERVER_HOST is not set")
	}

	connString := settings.ReadDbConnectionConfig()
	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Fatal("Unable to connect to DB", zap.Error(err))
	}
	defer connPool.Close()

	if err := connPool.Ping(ctx); err != nil {
		logger.Fatal("Unable to ping DB", zap.Error(err))
	}

	logger.Info("db connection success")

	uow := repository.NewUOW(connPool)

	eventRepo := repository.NewEventRepo(connPool)
	userRepo := repository.NewUserRepo(connPool)
	imageRepo := repository.NewImageRepo(connPool)

	eventSvc := service.NewEventService(&uow, &eventRepo, &userRepo, &imageRepo)
	jwtMgr := jwt.NewJwtManager(cfg.JWT.Secret)
	userSvc := service.NewUserService(&userRepo, &imageRepo, &eventRepo, &jwtMgr)
	searchEngine := service.NewSearchEngine(&eventRepo)
	validator := handlers.NewValidator()

	server := controller.NewServer(ctx, logger, &eventSvc, &userSvc, &searchEngine, &validator, &jwtMgr)
	if err != nil {
		logger.Fatal("error on create server: %v", zap.Error(err))
	}

	sc := echo.StartConfig{
		Address:         host + ":" + port,
		GracefulTimeout: 5 * time.Second,
	}

	fmt.Println("server started on http://" + host + ":" + port)
	fmt.Println("swagger ui on http://" + host + ":" + port + "/swagger/index.html")
	if err := sc.Start(ctx, server); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}

	logger.Info("server stop")
}
