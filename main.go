package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

type Action int

const (
	Start Action = iota
	Stop  Action = iota
	Reset Action = iota
)

var layout = "06-02-01 15:04:05"

func startMusic(app string) {
	log.Infof("[%v] Starting music on %s", time.Now().Format(layout), app)
	cmd := exec.Command("osascript", fmt.Sprintf("osascript -e 'tell app %s to play'", app))
	err := cmd.Run()
	if err != nil {
		log.Errorf("%v", err)
	}
}

func stopMusic(app string) {
	log.Infof("[%v] Stopping music on %s", time.Now().Format(layout), app)
	cmd := exec.Command("osascript", fmt.Sprintf("osascript -e 'tell app %s to pause'", app))
	err := cmd.Run()
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
