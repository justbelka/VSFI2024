#!/bin/bash

docker build -t workshop-0 .
docker run --rm -p 8080:80 --name workshop-0 workshop-0