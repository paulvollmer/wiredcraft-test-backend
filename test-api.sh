#!/bin/sh

rm -f server.db
./cmd/wiredcraft-test-backend-server/wiredcraft-test-backend-server --port 8000 &
sleep 2
newman run postman.json
pkill wiredcraft-test-backend-server
