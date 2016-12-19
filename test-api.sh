#!/bin/sh

rm -f server.db

./cmd/wiredcraft-test-backend-server/wiredcraft-test-backend-server --port 8000  & SERVER_PID=$!
sleep 2
newman run postman.json
kill -9 $SERVER_PID
