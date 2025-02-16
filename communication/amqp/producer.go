package amqp

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"log/slog"
	"tg-dispatcher/domain"
	"tg-dispatcher/lib/logger/sl"

	"github.com/rabbitmq/amqp091-go"
)

// Producer отвечает за отправку сообщений в RabbitMQ через Exchange
type Producer struct {
	conn       *amqp091.Connection
	channel    *amqp091.Channel
	exchange   string
	routingKey string
	log        *slog.Logger
}

// NewProducer создает нового продюсера и подключается к RabbitMQ
func NewProducer(amqpURL, exchange, routingKey string, log *slog.Logger) (*Producer, error) {
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

	// Объявляем Exchange (если он не объявлен заранее)
	err = ch.ExchangeDeclare(
		exchange, // Имя Exchange
		"direct", // Тип (direct, topic, fanout, headers)
		true,     // Долговечный (durable)
		false,    // Автоудаляемый (auto-delete)
		false,    // Внутренний
		false,    // Без ожидания подтверждения
		nil,      // Аргументы
	)
	if err != nil {
		conn.Close()
		ch.Close()
		return nil, err
	}

	return &Producer{
		conn:       conn,
		channel:    ch,
		exchange:   exchange,
		routingKey: routingKey,
		log:        log,
	}, nil
}

// StartPublishing читает сообщения из канала и отправляет их в RabbitMQ через Exchange
func (p *Producer) StartPublishing(messageChannel chan domain.MessageReq) {
	for msg := range messageChannel {
		err := p.PublishMessage(msg)
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
	}
}

// PublishMessage отправляет сообщение в RabbitMQ через Exchange
func (p *Producer) PublishMessage(msg domain.MessageReq) error {
	ctx, span := otel.Tracer("tg-dispatcher").Start(context.Background(), "PublishMessage")
	defer span.End()

	body, err := json.Marshal(msg.Message)
	if err != nil {
		return err
	}

	routingKey := msg.Destination.String()

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange, // Exchange
		routingKey, // Routing Key (для direct или topic exchange)
		false,      // Mandatory
		false,      // Immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers: amqp091.Table{
				"uuid": msg.UUID,
			},
		},
	)

	if err != nil {
		span.RecordError(err)
		p.log.Error("Ошибка публикации", sl.Err(err))
	} else {
		span.SetAttributes(attribute.String("uuid", msg.UUID))
		p.log.Info("Сообщение отправлено", slog.String("uuid", msg.UUID))
	}
	return err
}

// Close закрывает соединение с RabbitMQ
func (p *Producer) Close() {
	p.channel.Close()
	p.conn.Close()
}
