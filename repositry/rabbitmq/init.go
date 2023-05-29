package rabbitmq

import (
	"fmt"
	"mall/config"

	"github.com/streadway/amqp"
)

func NewRabbitMQConn(c *config.RabbitMQConf) (conn *amqp.Connection, err error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.User, c.Password, c.Host, c.Port)
	conn, err = amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}
	return
}
