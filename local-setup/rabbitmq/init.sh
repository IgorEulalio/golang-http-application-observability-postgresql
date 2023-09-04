#!/bin/bash

# Wait for RabbitMQ server to start
sleep 10

# Create 'repositories' queue
rabbitmqadmin -u guest -p guest declare queue name=repositories durable=true

# Run RabbitMQ server
docker-entrypoint.sh rabbitmq-server
