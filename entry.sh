#!/bin/bash
echo $1
if [ "$1" = "up" ]
then
	docker compose build && docker compose up && goose postgres "user=asymptoter password=password dbname=practice sslmode=disable" up
elif [ "$1" = "down" ]
then
	docker compose down
fi
