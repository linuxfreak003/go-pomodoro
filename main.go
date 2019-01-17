//go:generate protoc --go_out=plugins=grpc:./ pb/pomodoro.proto
package main

import "github.com/linuxfreak003/go-pomodoro/cmd"

func main() {
	cmd.Execute()
}
