package main

import (
	"log/slog"
	"os"
	"tg-dispatcher/communication/amqp"
	"tg-dispatcher/config"
	"tg-dispatcher/domain"
	"tg-dispatcher/lib/logger/sl"
	"tg-dispatcher/service"
	"tg-dispatcher/service/processors"
	"tg-dispatcher/storage/postgresql"
	"tg-dispatcher/tracing"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	shutdown := tracing.InitTracer(&cfg.OtlpConfig)
	defer shutdown()
	// Создаем канал для передачи сообщений
	inputMessageChannel := make(chan domain.Update, 100)
	outputMessageChannel := make(chan domain.MessageReq, 100)

	consumer := registerConsumer(inputMessageChannel, &cfg.AmqpConf, log)
	producer := registerProducer(outputMessageChannel, &cfg.AmqpConf, log)
	defer consumer.Close()
	defer producer.Close()

	storage, err := postgresql.New(cfg.PgConf.GetDbUri())
	if err != nil {
		log.Error("Cannot init db", sl.Err(err))
	}
	allProcessors := processors.CreateAllProcessors(outputMessageChannel, storage, storage, log)
	messageService := service.NewMessageProcessService(inputMessageChannel, outputMessageChannel, allProcessors, log)
	workerCount := 5
	messageService.StartProcessing(workerCount)

	// Держим программу живой
	for {
		time.Sleep(time.Second)
	}
}

func registerConsumer(inputMessageChannel chan domain.Update, cfg *config.AmqpConfig, log *slog.Logger) *amqp.Consumer {
	consumer, err := amqp.NewConsumer(cfg.GetAmqpUri(), cfg.QueueName, log)
	if err != nil {
		log.Error("Ошибка создания потребителя:", sl.Err(err))
	}

	go consumer.StartListening(inputMessageChannel)
	return consumer
}

func registerProducer(outputMessageChannel chan domain.MessageReq, cfg *config.AmqpConfig, log *slog.Logger) *amqp.Producer {
	producer, err := amqp.NewProducer(cfg.GetAmqpUri(), cfg.ExchangeName, cfg.RoutingKey, log)
	if err != nil {
		log.Error("Ошибка создания Producer:", sl.Err(err))
	}
	go producer.StartPublishing(outputMessageChannel)
	return producer
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
