#!/bin/bash
echo "Waiting for Kafka to start..."
sleep 30  # Даем Kafka время на запуск

kafka-topics.sh --bootstrap-server localhost:9092 --create --topic jwtCheckAnswer --partitions 4 --replication-factor 1
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic jwtCheckRequest --partitions 4 --replication-factor 1
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic userSignUp --partitions 2 --replication-factor 1
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic seasons --partitions 2 --replication-factor 1
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic dailyTasks --partitions 2 --replication-factor 1
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic gameSettings --partitions 2 --replication-factor 1

echo "Topics created successfully"