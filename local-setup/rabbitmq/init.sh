#!/bin/bash

# Wait for RabbitMQ server to start
sleep 10

# Check if the 'repositories' queue already exists
if ! rabbitmqadmin -u guest -p guest list queues name | grep -q 'repositories'; then
  # Create 'repositories' queue
  rabbitmqadmin -u guest -p guest declare queue name=repositories durable=true
fi

# Run RabbitMQ server
docker-entrypoint.sh rabbitmq-server

