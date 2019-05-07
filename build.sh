#!/bin/sh

cp conf/config.json bin/config.json
go build -o bin/microservice .
docker build -t $1/$2:$3 .
