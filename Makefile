all: run

build:
	go build

generate:
	protoc --go_out=./ pb/pomodoro.proto

run: build
	./go-pomodoro
