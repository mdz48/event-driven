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
	orderQueue    string
	orderExchange string
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	// Si la URL está vacía, usa un valor predeterminado
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error abriendo canal: %s", err)
	}

	// Definir las constantes para exchanges y colas que ya existen
	orderExchange := "orders.events"

	// Ya no declaramos el exchange ni las colas - asumimos que existen
	log.Println("Usando configuración manual de RabbitMQ - exchange y colas ya existentes")

	return &RabbitMQ{
		conn:          conn,
		channel:       ch,
		orderQueue:    "orders.created", // Cola predeterminada para el consumidor de ejemplo
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
		eventType,       // routing key (ahora es sólo el tipo de evento)
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

// Consumir mensajes para un tipo de evento específico
func (r *RabbitMQ) ConsumeSpecificEvents(eventType string, handler func(order domain.Order) error) error {
	queueName := "orders." + eventType
	
	// Ya no declaramos la cola ni los bindings - asumimos que existen
	
	msgs, err := r.channel.Consume(
		queueName, // nombre de la cola específica del evento
		"",        // consumidor
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("error registrando consumidor para %s: %s", eventType, err)
	}

	// Resto del código igual...
	go func() {
		for d := range msgs {
			var order domain.Order
			err := json.Unmarshal(d.Body, &order)

			if err != nil {
				log.Printf("Error decodificando mensaje: %s", err)
				d.Nack(false, true) // rechazar el mensaje y ponerlo de nuevo en la cola
				continue
			}

			err = handler(order)
			if err != nil {
				log.Printf("Error manejando evento %s: %s", eventType, err)
				d.Nack(false, true) // rechazar y reintentar
			} else {
				d.Ack(false) // confirmar que se procesó correctamente
			}
		}
	}()

	log.Printf("Consumidor de eventos '%s' iniciado en cola existente", eventType)
	return nil
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
