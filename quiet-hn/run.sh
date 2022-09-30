#!/bin/sh

# Install
# go install github.com/cosmtrek/air@latest

# go run main.go \
# ./main \
air -- \
    --port 8080 \
    --n 0 \
    --caching "mem" \
    --cache-timeout 60m \
    --redis-host "localhost" \
    --redis-port 6379 \
    --redis-password="" \
    --rate-type "none" \
    --rate-burst 5 \
    --rate-interval 3s \
    --debug
