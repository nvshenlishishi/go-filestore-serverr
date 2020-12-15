## 9-1 Ubuntu下通过Docker按照RabbitMQ

## 9-2 关于任务的同步与异步
### 同步与异步

###  异步逻辑架构


## 9-3 RabbitMQ简介
- RabbitMQ是什么
``` 
一种开源的消息代码
一种面向消息的中间件
一种默认遵循AMQP协议的MQ服务
```
- RabbitMQ可以解决什么
``` 
逻辑解耦
异步任务
消息持久化
重启不影响
削峰，大规模消息处理
```
- RabbitMQ的特点
``` 
可靠性     持久化、传输确认、发布确认
可扩展性    多个节点可以组成一个集群，可动态更改
多语言客户端 几乎支持所有常用语言
管理界面    易用的用户界面，便于监控和管理
```
## 9-4 RabbitMQ工作原来和转发模式
### RabbitMQ工作原理

### RabbitMQ关键术语
- Exchange: 消息交换机，决定消息按什么规则，路由到哪个队列
- Queue: 消息载体，每个消息都会被投到一个或者多个队列
- Binding:绑定，把exchange和queue按照路由规则绑定起来
- Routing Key:路由关键字,exchange根据这关键字来投递消息
- Channel: 消息通道，客户端的每个连接建立多个Channel
- Producer: 消息生产者，用于投递消息的程序
- Consumer: 消息消费者，用于接受消息的程序

### Exchange工作模式
- Fanout: 类似广播，转发到所有绑定交换机的Queue
- Direct:类似单播，RoutingKey 和 BindingKey完全匹配
- Topic: 类似组播，转发到符合通配符的Queue
- Headers:请求头与消息头匹配，才能接收消息

## 9-5 docker安装rabbitMq及UI管理
``` 
mkdir -p /data/rabbitmq 
docker run -d --hostname rabbit-svr --name rabbit -p 5672:5672 -p 15672:15672 
-p 25672:25672 
-v /data/rabbitmq:/var/lib/rabbitmq 
rabbitmq:management
```
- 登录rabbitmq
``` 
localhost:15672
guest
guest
```
- 创建exchanges
``` 
# add a new exchange 
Name: uploadserver.trans
Type: direct
Durability: Durable
```
- 创建Queues
``` 
# add a new queue
Name: uploadserver.trans.oss
Durability: Durable
```
- 添加绑定
``` 
# add binding to this queue
From exchange: uploadserver.trans
Routing Key:oss
```
- 发布信息
``` 
# publish message
Routing key: oss
Payload:
test
```
- 获取信息
``` 
# get messages
Ack Mode: Nack message requeue true
```
## 9-6 实现异步转移的mq生产者
- rabbitmq.go
- producer.go
## 9-7 实现异步转移的mq消费者
- consumer.go

## 9-8 异步转移文件测试+小结
### 小结
- 同步与异步逻辑的对比分析
- RabbitMq的基本概念与工作原理
- RabbitMq的安装与控制台管理
- Go结合队列实现文件的异步转移
- 文件上传与转移的简单测试