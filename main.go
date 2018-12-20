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

func startMusic(app string) {
	log.Infof("[%v] Starting music on %s", time.Now().Format(timeLayout), app)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = dbus(app, "Play")
	case "darwin":
		err = osascript(app, "play")
	default:
		log.Fatalf("unknown operating system")
	}
	if err != nil {
		log.Errorf("%v", err)
	}
}

func stopMusic(app string) {
	log.Infof("[%v] Stopping music on %s", time.Now().Format(timeLayout), app)
	var err error
	switch runtime.GOOS {
	case "linux":
		err = dbus(app, "Pause")
	case "darwin":
		err = osascript(app, "pause")
	default:
		log.Fatalf("unknown operating system")
	}
	if err != nil {
		log.Errorf("%v", err)
	}
}

func Timer(actions chan Action, app string, minutes, interval int) {
	startMusic(app)
	timer1 := time.NewTimer(time.Duration(minutes) * time.Minute)
	log.Infof("Timer set for %d minutes", minutes)
	timer2 := &time.Timer{
		C: make(chan time.Time),
	}
	for {
		select {
		case <-timer1.C:
			stopMusic(app)
			timer2 = time.NewTimer(time.Duration(interval) * time.Minute)
			log.Infof("Timer set for %d minutes", interval)
		case <-timer2.C:
			startMusic(app)
			timer1 = time.NewTimer(time.Duration(minutes) * time.Minute)
			log.Infof("Timer set for %d minutes", minutes)
		case a := <-actions:
			switch a {
			case Start:
				startMusic(app)
			case Stop:
				stopMusic(app)
			case Reset:
				startMusic(app)
				timer1 = time.NewTimer(time.Duration(minutes) * time.Minute)
				log.Infof("Timer set for %d minutes", minutes)
			}
		}
	}
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println(runtime.GOOS)
	}
	var minutes, interval int
	var app string
	flag.IntVar(&minutes, "length", 25, "pomodoro length in minutes")
	flag.IntVar(&interval, "break", 5, "pomodoro break in minutes")
	flag.StringVar(&app, "app", "spotify", "music app to use")
	flag.Parse()

	done := make(chan struct{})
	actions := make(chan Action)

	go Timer(actions, app, minutes, interval)

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
