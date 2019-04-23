#!/bin/bash
cd ../phenix
CGO_ENABLED=0 GOOS=linux go build -o node
mv node ../docker/
cd ../docker
docker build -t phenix-server:v1.0 .
rm node