#!/bin/bash

docker build -t workshop-1-small -f Dockerfile-small .
docker run --rm --name workshop-1-small workshop-1-small