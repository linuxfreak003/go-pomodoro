all: build

build: generate
	go build -o go-pomodoro ./cli

generate: pb/pomodoro.proto
	protoc --go_out=plugins=grpc:./ pb/pomodoro.proto

run: build
	./go-pomodoro
