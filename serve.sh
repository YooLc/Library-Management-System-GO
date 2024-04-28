#!/bin/bash

mkdir -p log

if [ "$1" = "start" ]; then
    echo "Starting MySQL server..."
    sudo systemctl start mysqld

    echo "Starting backend server..."
    go build .
    nohup go run main.go > log/backend.log &
    echo $! > log/backend.pid

    echo "Starting frontend server..."
    cd frontend
    pnpm install
    pnpm run build
    nohup pnpm run serve > ../log/frontend.log &
    echo $! > ../log/frontend.pid

    echo "Server started!"
elif [ "$1" = "stop" ]; then
    echo "Stopping backend server..."
    kill -9 $(cat log/backend.pid)
    rm log/backend.pid

    echo "Stopping frontend server..."
    kill -9 $(cat log/frontend.pid)
    rm log/frontend.pid

    echo "Server stopped!"
else
    echo "Invalid argument. Please use 'start' or 'stop'."
fi