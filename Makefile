all: run

build: generate
	go build

generate:
	protoc --go_out=./ pb/pomodoro.proto

run: build
	./go-pomodoro
