SPEC=spec.yml

all: gen install build test docs release

validate:
	@echo "==> Validate"
	@swagger validate ${SPEC}

gen: validate
	@echo "==> Generate"
	@swagger generate server --spec ${SPEC}

install:
	@go get -u -f -v ./...

build:
	@echo "==> Build"
	@cd cmd/wiredcraft-test-backend-server && go build -v

test:
	@echo "==> Test"
	@cd restapi && go test -v
	@rm -f restapi/*.db

test-api:
	./cmd/wiredcraft-test-backend-server/wiredcraft-test-backend-server --port 8000 & FUZZ_PID=$! && sleep 2 && newman run postman.json && kill -9 ${FUZZ_PID}


serve:
	./cmd/wiredcraft-test-backend-server/wiredcraft-test-backend-server --port 8000

docs:
	@echo "==> Docs"
	# TODO: create the docs

release: build
	@echo "==> Release"
	# TODO: create a release package

clean:
	@rm -rf cmd models
	@rm -rf restapi/operations
	@rm -f restapi/doc.go
	@rm -f restapi/embedded_spec.go
	@rm -f restapi/server.go
	@rm -f *.log
	@rm -f *.db
