all: run

build:
	go build -o go-pomodoro ./cli

generate:
	protoc --go_out=plugins=grpc:./ pb/pomodoro.proto

run: build
	./go-pomodoro
