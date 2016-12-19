# Wiredcraft Back-end Developer Coding Test

## Background

Build a restful api that could `CRUD (create, read, update, delete)` user data from a persistence database.

## Architecture

The Restful API is specified in an openAPI aka swagger format.  
The server is build on top of `go-swagger` and `boltDB` (Embedded Database) to store the Data.


## How to run the code
*Requirements*
- golang to compile the final API server
- nodejs (for API tests with newman, the [postman](https://www.getpostman.com/) cli tool)

Clone the Repository and simple run
```
cd $GOPATH/src/github.com/paulvollmer
git clone https://github.com/paulvollmer/wiredcraft-test-backend
cd wiredcraft-test-backend

# The make Task generate and build the server
make
```

## API Documentation
The API Docs can be found at the `http://localhost:8000/docs`

## Unit Test
There are two kind of tests.
- Database tests
```
make test
```
- Server tests
```
# This test requires a running server to send the test requests
# start running a server...
make serve
# execute the test collection
make test-api
```
