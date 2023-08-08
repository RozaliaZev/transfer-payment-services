# transfer-payment-services

## Схема работы

Есть два сервиса, общающиеся между собой через брокер сообщений.

![alt text](https://sun9-53.userapi.com/impg/f18Z4Ad_oMDLJBpsHUl6By5qvfK-bWipXz-9ww/_f4PplVy2ZQ.jpg?size=964x152&quality=96&sign=59f8525751c467cee2847f9c73c97574&type=album)

## Конфигурирование сервиса

Запуск брокера через docker-compose

- **db** информация для подключения к базе данных PostgreSQL

```sh
version: "3.9"
services:
  kafka:
    image: wurstmeister/kafka
    hostname: kafka
    ports:
      - "9092:9092"
    links:
      - zookeeper
    environment:
      KAFKA_ADVERTISED_HOST_NAME: "127.0.0.1"
      KAFKA_ADVERTISED_PORT: "9092"
      KAFKA_ZOOKEEPER_CONNECT: "zookeeper:2181"
      KAFKA_CREATE_TOPICS: |
        test-topic:1:1,
        transfer_request_topic:1:1,
        check_result_topic:1:1,
        check_status_topic:1:1
      KAFKA_BROKER: "localhost:9092"
    depends_on:
      - zookeeper
    container_name: tgbot-kafka
  zookeeper:
      image: wurstmeister/zookeeper
      ports:
        - "2181:2181"
      container_name: tgbot-zookeeper
```

## Запросы

- **POST** загрузка новой транзакции. 

```sh
/transaction
```

```sh
{
  "senderId": "string",
  "requestId": "string",
  "amount": 0
}
```

В запросе передаётся сумма и id заявки, реализована защиту - нельзя создавать две заявки с одинковым id у одного пользователя, поэтому в запросе присутствует id пользователя.
