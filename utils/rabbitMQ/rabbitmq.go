package rabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
)

func InitMqConn() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("connect to RabbitMQ failed, err:%v\n", err)
		panic(err)
	}
	return conn, nil
}
func ConnClose(conn *amqp.Connection) {
	conn.Close()
}
