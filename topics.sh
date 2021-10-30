#!/bin/bash

sleep 15 && kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 3 --topic quick-starter)&) && /etc/confluent/docker/run 