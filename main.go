package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

type Action int

const (
	Start Action = iota
	Stop  Action = iota
)

func startMusic(app string, minutes int) {
	fmt.Printf("Starting music on %s\n", app)
	/*
		cmd := exec.Command("osascript", fmt.Sprintf("osascript -e 'tell app %s to play'", app))
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	*/
}

func Timer(actions chan Action, app string, minutes int) {
	for {
		switch <-actions {
		case Start:
			startMusic(app, minutes)
		case Stop:
			fmt.Println("Stop action received")
		}
	}
}

func main() {
	var minutes, interval int
	var app string
	flag.IntVar(&minutes, "length", 25, "pomodoro length in minutes")
	flag.IntVar(&interval, "break", 25, "pomodoro break in minutes")
	flag.StringVar(&app, "app", "spotify", "music app to use")
	flag.Parse()

	done := make(chan struct{})
	actions := make(chan Action)

	go Timer(actions, app, minutes)

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
