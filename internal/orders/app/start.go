package app

import (
	"kaf-interface/internal/orders/config"
	"kaf-interface/internal/orders/consumer"
	httpHandler "kaf-interface/internal/orders/handler/http_handler"
	kafkaHandler "kaf-interface/internal/orders/handler/kafka_handler"
	"kaf-interface/internal/orders/repository/cache"
	"kaf-interface/internal/orders/repository/postgres"
	"kaf-interface/internal/orders/server"
	"kaf-interface/internal/orders/service"
	"kaf-interface/pkg/cacheMap"
	"kaf-interface/pkg/logger"
	"kaf-interface/pkg/postgresDB"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func StartApp() {

	// загрузка конфига
	config, err := config.MustLoad("config")
	if err != nil {
		log.Fatal(err)
		return
	}

	logger := logger.Load(config) // загрузка логгера

	// инициализация SQL storage
	postgresDB, err := postgresDB.NewPostgres(config)
	if err != nil {
		logger.Error("load db error", slog.Any("error", err))
		return
	}

	// создание Cache storage
	cacheMap := cacheMap.NewCacheMap()

	// инициализация слоя репозиториев
	postgresRepo := postgres.NewRepository(postgresDB)
	cacheRepo := cache.NewRepository(cacheMap)

	deps := service.Deps{
		DBRepo:    postgresRepo,
		CacheRepo: cacheRepo,
		Config:    config,
	}

	// инициализация слоя сервиса
	service := service.NewService(deps)

	// создание консьюмера
	kafkaHandler := kafkaHandler.NewHandler(service, logger)
	c, err := consumer.NewConsumer(kafkaHandler, logger, config)
	if err != nil {
		logger.Error("error with new consumer", slog.Any("error", err))
		return
	}

	go c.Start() // запуск консьюмера

	// создание WEB хендлера
	httpHandler := httpHandler.NewHandler(service, logger).InitRouters(config)
	server := server.NewServer(config, httpHandler) // создание сервера
	if err := server.Run(); err != nil {            // запуск сервера
		logger.Error("error with start server", slog.Any("error", err))
		return
	}

	// gracefull shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Error("stop error", slog.Any("error", c.Stop()))

}
