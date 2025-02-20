package main

import (
    "fmt"
    "github.com/edebernis/sizematch-protobuf/go/items"
    "github.com/golang/protobuf/proto"
    "github.com/streadway/amqp"
    "time"
)

type messenger struct {
    host       string
    port       string
    username   string
    password   string
    vhost      string
    appID      string
    connection *amqp.Connection
    channel    *amqp.Channel
}

func (m *messenger) buildURL() string {
    return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", m.username, m.password, m.host, m.port, m.vhost)
}

func (m *messenger) connect(connectionAttempts int) error {
    var err error
    url := m.buildURL()

    m.connection, err = amqp.Dial(url)
    if err != nil {
        if connectionAttempts < 1 {
            return err
        }
        time.Sleep(5 * time.Second)
        return m.connect(connectionAttempts - 1)
    }

    m.channel, err = m.connection.Channel()
    if err != nil {
        return err
    }

    return nil
}

func (m *messenger) setupPublisher(exchangeName, routingKey, queueName string) error {
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

func (m *messenger) setupConsumer(queueName string, prefetchCount int) error {
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

func (m *messenger) publishItem(exchangeName, routingKey string, item *items.NormalizedItem) error {
    body, err := proto.Marshal(item)
    if err != nil {
        return err
    }

    msg := amqp.Publishing{
        ContentType: "application/protobuf",
        AppId:       m.appID,
        Body:        body,
    }

    err = m.channel.Publish(exchangeName, routingKey, true, false, msg)
    if err != nil {
        return err
    }

    return nil
}

func (m *messenger) consumeItem(queueName string, callback func(item *items.Item) error) error {
    msgs, err := m.channel.Consume(queueName, "", false, false, false, false, nil)
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            item := items.Item{}
            err := proto.Unmarshal(msg.Body, &item)
            if err != nil {
                fmt.Println("could not decode protobuf item: " + err.Error())
                msg.Nack(false, false)
                continue
            }

            err = callback(&item)
            if err != nil {
                msg.Nack(false, false)
                continue
            }

            msg.Ack(false)
        }
    }()

    return nil
}

func (m *messenger) close() {
    m.connection.Close()
}
