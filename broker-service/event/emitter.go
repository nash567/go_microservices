package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func (e Emitter) setup() error {
	channel, err := e.conn.Channel()

	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.conn.Channel()

	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel...")
	err = channel.Publish(
		"log_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		conn: conn,
	}
	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
