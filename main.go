package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Action int

const (
	Start Action = iota
	Stop  Action = iota
	Reset Action = iota
)

func startMusic(app string) {
	fmt.Printf("Starting music on %s\n", app)
	cmd := exec.Command("osascript", fmt.Sprintf("osascript -e 'tell app %s to play'", app))
	_ = cmd.Run()
}

func stopMusic(app string) {
	fmt.Printf("Stopping music on %s\n", app)
	cmd := exec.Command("osascript", fmt.Sprintf("osascript -e 'tell app %s to pause'", app))
	_ = cmd.Run()
}

func Timer(actions chan Action, app string, minutes, interval int) {
	timer1 := time.NewTimer(time.Duration(minutes) * time.Minute)
	startMusic(app)
	timer2 := &time.Timer{
		C: make(chan time.Time),
	}
	for {
		select {
		case <-timer1.C:
			stopMusic(app)
			timer2 = time.NewTimer(time.Duration(interval) * time.Minute)
		case <-timer2.C:
			startMusic(app)
			timer1 = time.NewTimer(time.Duration(minutes) * time.Minute)
		case a := <-actions:
			switch a {
			case Start:
				startMusic(app)
			case Stop:
				stopMusic(app)
			case Reset:
				timer1 = time.NewTimer(time.Duration(minutes) * time.Minute)
				startMusic(app)
			}
		}
	}
}

func main() {
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
			fmt.Printf("> ")
			scanner.Scan()
			text := scanner.Text()
			switch text {
			case "start":
				actions <- Start
			case "stop":
				actions <- Stop
			case "q":
				done <- struct{}{}
			default:
				fmt.Printf("command not recognized: %s\n", text)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	<-done
}
