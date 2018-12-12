all: run

build:
	go build

run: build
	./go-pomodoro
