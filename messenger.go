package main

import (
    "fmt"
    "github.com/streadway/amqp"
)

// Messenger consumes and publishes msgs to RabbitMQ
type Messenger struct {
    host       string
    port       int
    username   string
    password   string
    vhost      string
    appID      string
    connection *amqp.Connection
    channel    *amqp.Channel
}

func (m *Messenger) buildURL() string {
    return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", m.username, m.password, m.host, m.port, m.vhost)
}

// Connect to RabbitMQ
func (m *Messenger) Connect() error {
    var err error

    url := m.buildURL()
    m.connection, err = amqp.Dial(url)
    if err != nil {
        return err
    }

    m.channel, err = m.connection.Channel()
    if err != nil {
        return err
    }

    return nil
}

// SetupPublisher setups exchange and queue for publishing
func (m *Messenger) SetupPublisher(exchangeName, routingKey, queueName string) error {
    err := m.channel.ExchangeDeclare(exchangeName, "direct", false, false, false, false, nil)
    if err != nil {
        return err
    }

    _, err = m.channel.QueueDeclare(queueName, false, false, false, false, nil)
    if err != nil {
        return err
    }

    err = m.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
    if err != nil {
        return err
    }

    return nil
}

// SetupConsumer setups queue for consuming
func (m *Messenger) SetupConsumer(queueName string, prefetchCount int) error {
    _, err := m.channel.QueueDeclare(queueName, false, false, false, false, nil)
    if err != nil {
        return err
    }

    err = m.channel.Qos(prefetchCount, 0, false)
    if err != nil {
        return err
    }

    return nil
}

// Publish messages
func (m *Messenger) Publish(exchangeName, routingKey string, body []byte) error {
    msg := amqp.Publishing{ContentType: "application/json", AppId: m.appID, Body: body}

    err := m.channel.Publish(exchangeName, routingKey, true, false, msg)
    if err != nil {
        return err
    }

    return nil
}

// Consume messages
func (m *Messenger) Consume(queueName string, callback func([]byte) error) error {
    msgs, err := m.channel.Consume(queueName, "", false, false, false, false, nil)
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            err = callback(msg.Body)
            if err != nil {
                msg.Nack(false, true)
            }
            msg.Ack(false)
        }
    }()

    return nil
}

// Close connection to RabbitMQ
func (m *Messenger) Close() {
    m.connection.Close()
}
