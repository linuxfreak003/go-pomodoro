# go-pomodoro
Pomodoro timer written in Go

## Install

`go get github.com/linuxfreak003/go-pomodoro`

OR

1. Clone the repo
1. `cd` into the repo
1. run `go install`

## Run

go-pomodoro is setup with a server and client model. To run you will first need to run `go-pomodoro server` or `go-pomodoro client`

### Server

The possible options on the server are:
* `-p <port>` - default 50051
* `-t <slack_token>` - default ""
* `-c <slack_channel` - default "pomodoro-spotify"

### Client

The possible options on the client are:
* `--app <application>` - default "spotify" (spotify is also the only app that has been tested)
* `--host <host_ip>` - default "127.0.0.1"
* `-p <port>` - default 50051

The options will also be shown if you add `-h` to any of the commands

When the client is open there are some commands you can do
* `p` previous track (only linux)
* `n` next track (only linux)
* `<space>` Pause/Play
* `q` quit
