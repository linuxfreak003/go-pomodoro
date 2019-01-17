all: build

build: generate
	go build

generate: pb/pomodoro.proto
	go generate

run: build
	./go-pomodoro

clean:
	rm -f go-pomodoro
	rm -f pb/*.pb.go
