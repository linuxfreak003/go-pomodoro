package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Action int

const (
	Start Action = iota
	Stop  Action = iota
	Reset Action = iota
)

var timeLayout = "06-02-01 15:04:05"

func dbus(app, cs string) error {
	cmdString := fmt.Sprintf("--print-reply --dest=org.mpris.MediaPlayer2.%s /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.%s", app, cs)
	cmd := exec.Command("dbus-send", strings.Split(cmdString, " ")...)
	return cmd.Run()
}
func osascript(app, cs string) error {
	cmdString := fmt.Sprintf("tell app \"%s\" to %s", app, cs)
	cmd := exec.Command("osascript", "-e", cmdString)
	return cmd.Run()
}
func musicCommand(app, command string) error {
	log.Infof("[%v] %s music on %s", time.Now().Format(timeLayout), command, app)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = dbus(app, command)
	case "darwin":
		err = osascript(app, command)
	default:
		err = fmt.Errorf("unknown or unsupported operating system")
	}
	if err != nil {
		log.Errorf("%v", err)
	}
	return err
}

func StartTimer(minutes int, app, command string) *time.Timer {
	musicCommand(app, command)
	log.Infof("Timer set for %d minutes", minutes)
	return time.NewTimer(time.Duration(minutes) * time.Minute)
}

func PomodoroTimer(actions chan Action, app string, start, minutes, interval int) {
	timer1 := StartTimer(start, app, "Play")
	timer2 := &time.Timer{C: make(chan time.Time)}
	for {
		select {
		case <-timer1.C:
			timer2 = StartTimer(interval, app, "Pause")
		case <-timer2.C:
			timer1 = StartTimer(minutes, app, "Play")
		case a := <-actions:
			switch a {
			case Start:
				musicCommand(app, "Play")
			case Stop:
				musicCommand(app, "Pause")
			case Reset:
				timer1 = StartTimer(minutes, app, "Play")
			}
		}
	}
}

func main() {
	var minutes, interval, start int
	var app string
	flag.IntVar(&start, "start", 25, "starting point for time")
	flag.IntVar(&minutes, "length", 25, "pomodoro length in minutes")
	flag.IntVar(&interval, "break", 5, "pomodoro break in minutes")
	flag.StringVar(&app, "app", "spotify", "music app to use")
	flag.Parse()

	done := make(chan struct{})
	actions := make(chan Action)

	go PomodoroTimer(actions, app, start, minutes, interval)

	scanner := bufio.NewScanner(os.Stdin)

	// command line interface
	go func() {
		for {
			time.Sleep(time.Second * 1)
			fmt.Printf("$ ")
			scanner.Scan()
			text := scanner.Text()
			switch text {
			case "start":
				actions <- Start
			case "stop":
				actions <- Stop
			case "reset":
				actions <- Reset
			case "q", "quit", "exit":
				done <- struct{}{}
			default:
				log.Warnf("command not recognized: %s\n", text)
			}
		}
	}()

	<-done
}
