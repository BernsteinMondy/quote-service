#!/bin/bash

if [ -f .env.dev ]; then
    export $(grep -v '^#' .env.dev | xargs)
else
    echo "Файл .env.dev не найден."
    exit 1
fi

goose -dir migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER dbname=$DB_NAME password=$DB_PASSWORD sslmode=$DB_SSL_MODE" up

cd ./src/cmd && go run .