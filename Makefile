MAIN = ./cmd/server/main.go
BIN = ./bin/fiber_jwt

all: build

build:
	go build -o $(BIN) $(MAIN)

dev:
	CompileDaemon -build="make" -command="$(BIN)"

clean:
	rm -r ./bin
