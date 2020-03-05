#!/bin/bash
docker-compose build && docker-compose up && goose postgres "user=asymptoter password=password dbname=practice sslmode=disable" up
