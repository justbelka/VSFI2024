#!/bin/bash

docker build -t workshop-1-big -f Dockerfile-big .
docker run --rm --name workshop-1-big workshop-1-big