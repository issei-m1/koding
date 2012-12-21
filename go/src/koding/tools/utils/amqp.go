package utils

import (
	"fmt"
	"github.com/streadway/amqp"
	"koding/config"
	"koding/tools/log"
	"strings"
	"time"
)

func AmqpAutoReconnect(component string, handler func(consumeConn, publishConn *amqp.Connection)) {
	for !ShuttingDown {
		func() {
			defer time.Sleep(time.Second)
			defer log.RecoverAndLog()

			log.Info("Connecting to AMQP server...")

			consumeConn := CreateAmqpConnection(component)
			defer consumeConn.Close()
			publishConn := CreateAmqpConnection(component)
			defer publishConn.Close()

			log.Info("Successfully connected to AMQP server.")

			handler(consumeConn, publishConn)

			if !ShuttingDown {
				log.Warn("Connection to AMQP server lost.")
			}
		}()
	}
}

func CreateAmqpConnection(component string) *amqp.Connection {
	user := config.Current.AmqpUser
	user = strings.Replace(user, "<component>", component, 1)
	conn, err := amqp.Dial("amqp://" + user + ":" + config.Current.AmqpPassword + "@" + config.Current.AmqpHost)
	if err != nil {
		panic(err)
	}
	go func() {
		for err := range conn.NotifyClose(make(chan *amqp.Error)) {
			log.Warn("AMQP connection: " + err.Error())
		}
	}()
	return conn
}

func CreateAmqpChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	go func() {
		for err := range channel.NotifyClose(make(chan *amqp.Error)) {
			log.Warn("AMQP channel: " + err.Error())
		}
	}()
	return channel
}

func DeclareBindConsumeAmqpQueue(channel *amqp.Channel, kind, exchange, key string) <-chan amqp.Delivery {
	if err := channel.ExchangeDeclare(exchange, kind, false, true, false, false, nil); err != nil {
		panic(err)
	}

	if _, err := channel.QueueDeclare("", false, true, false, false, nil); err != nil {
		panic(err)
	}

	if err := channel.QueueBind("", key, exchange, false, nil); err != nil {
		panic(err)
	}

	stream, err := channel.Consume("", "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	return stream
}

func DeclareAmqpPresenceExchange(channel *amqp.Channel, exchange, serviceType, serviceGenericName, serviceUniqueName string) {
	if err := channel.ExchangeDeclare(exchange, "x-presence", false, true, false, false, nil); err != nil {
		panic(err)
	}

	if _, err := channel.QueueDeclare("", false, true, true, false, nil); err != nil {
		panic(err)
	}

	routingKey := fmt.Sprintf("serviceType.%s.serviceGenericName.%s.serviceUniqueName.%s", serviceType, serviceGenericName, serviceUniqueName)
	if err := channel.QueueBind("", routingKey, exchange, false, nil); err != nil {
		panic(err)
	}
}
