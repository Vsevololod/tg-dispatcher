package amqp

import (
	"log/slog"
	"tg-dispatcher/domain"
	"tg-dispatcher/lib/logger/sl"

	"github.com/rabbitmq/amqp091-go"
)

// Consumer отвечает за подключение к RabbitMQ и отправку сообщений в канал
type Consumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   string
	log     *slog.Logger
}

// NewConsumer создает нового потребителя
func NewConsumer(amqpURL, queueName string, log *slog.Logger) (*Consumer, error) {
	log.Info("Create Consumer")
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		conn.Close()
		ch.Close()
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   queueName,
		log:     log,
	}, nil
}

// StartListening запускает прослушивание очереди и отправку сообщений в канал
func (c *Consumer) StartListening(messageChannel chan domain.Update) {
	const op = "Consumer.StartListening"
	log := c.log.With(
		slog.String("op", op),
	)
	log.Info("Start listening")
	msgs, err := c.channel.Consume(
		c.queue, "", false, false, false, false, nil,
	)
	if err != nil {
		log.Error("Ошибка подписки на очередь:", sl.Err(err))
	}

	// Обрабатываем каждое сообщение в горутине
	go func() {
		for msg := range msgs {
			update, err := domain.ParseUpdate(msg.Body)
			if value, ok := msg.Headers["uuid"].(string); ok {
				update.UUID = value
			}
			if err != nil {
				log.Error("Ошибка декодирования JSON: %s", sl.Err(err))
				_ = msg.Ack(false)
				continue
			}

			// Пишем сообщение в канал
			messageChannel <- update
			_ = msg.Ack(true) // Подтверждаем получение
		}
	}()
}

// Close закрывает соединение с RabbitMQ
func (c *Consumer) Close() {
	const op = "Consumer.Close"

	c.log.Info("Close consumer", slog.String("op", op))
	c.channel.Close()
	c.conn.Close()
}
