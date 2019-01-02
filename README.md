# go-pomodoro
pomodoro timer written in Go

## Install

`go get github.com/linuxfreak003/go-pomodoro`

## Run

Running `go-pomodoro` without any arguments will set the timer to a 25 minute focus and a 5 minute break.

Possible options for arguments:
 * `-start #` - Sets the initial starting point of the timer - default 25 min
 * `-length #` - Sets the length of the focus time - default 25
 * `-break #` - Sets the length of the break time - default 5
 * `-app <string>` - Specifies what app to use (only spotify tested)
