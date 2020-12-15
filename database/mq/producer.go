package mq

import (
	"github.com/streadway/amqp"
	"go-filestore-server/config"
	"go-filestore-server/logger"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

// UpdateRabbitHost : 更新mq host
func UpdateRabbitHost(host string) {
	config.DefaultConfig.RabbitURL = host
}

func InitMq() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.DefaultConfig.AsyncTransferEnable {
		return
	}
	if initChannel() {
		channel.NotifyClose(notifyClose)
	}

	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				logger.Infof("onNotifyChannelClosed:%+v\n", msg)
				initChannel()
			}
		}
	}()
}

func initChannel() bool {
	if channel != nil {
		return true
	}

	conn, err := amqp.Dial(config.DefaultConfig.RabbitURL)
	if err != nil {
		logger.Info(err.Error())
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	return true
}

// 发布消息
func Publish(exchange, routingKey string, msg []byte) bool {
	if !initChannel() {
		return false
	}

	if nil == channel.Publish(
		exchange,
		routingKey,
		false, // 如果没有对应的queue，就会丢弃这条消息
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		}) {
		return true
	}
	return false
}
