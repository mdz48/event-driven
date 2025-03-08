package producer

import (
	"context"
	"encoding/json"
	"event-driven/src/features/orders/domain"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	orderExchange string
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error abriendo canal: %s", err)
	}

	// Definir las constantes para exchanges
	orderExchange := "orders.events"

	// No declaramos nada - asumimos que existe
	log.Println("Usando configuración manual de RabbitMQ - exchange ya existente")

	return &RabbitMQ{
		conn:          conn,
		channel:       ch,
		orderExchange: orderExchange,
	}, nil
}

// Publicar eventos relacionados con órdenes
func (r *RabbitMQ) PublishOrderEvent(eventType string, order domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error codificando orden: %s", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		r.orderExchange, // exchange
		eventType,       // routing key (tipo de evento)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		})

	if err != nil {
		return fmt.Errorf("error publicando mensaje: %s", err)
	}

	log.Printf("Evento de orden publicado: %s, ID: %d", eventType, order.ID)
	return nil
}

// Ayudantes para diferentes eventos
func (r *RabbitMQ) NotifyOrderCreated(order domain.Order) error {
	return r.PublishOrderEvent("created", order)
}

func (r *RabbitMQ) NotifyOrderStatusChanged(order domain.Order) error {
	return r.PublishOrderEvent("status_changed", order)
}

func (r *RabbitMQ) NotifyOrderDeleted(orderID int) error {
	order := domain.Order{ID: orderID}
	return r.PublishOrderEvent("deleted", order)
}

// Cierra las conexiones de RabbitMQ
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
