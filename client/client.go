package client

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	pb "github.com/linuxfreak003/go-pomodoro/pb"
	"github.com/linuxfreak003/timer"
	term "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Action int

const (
	Play     Action = iota
	Pause    Action = iota
	Toggle   Action = iota
	Previous Action = iota
	Next     Action = iota

	Reset     Action = iota
	Remaining Action = iota
)

type Profile struct {
	Name string
	Host string
	Port uint16
}

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
	fmt.Println()
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

func Timer(actions chan Action, app string, profile Profile) {
	ctx := context.Background()
	addr := fmt.Sprintf("%s:%d", profile.Host, profile.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("%v", err)
	}
	defer conn.Close()

	client := pb.NewPomodoroClient(conn)
	stateTimer := &timer.Timer{}

	// Display remaining time
	go func() {
		for {
			time.Sleep(time.Second * 1)
			fmt.Printf("\r> %s", stateTimer.Remaining().Truncate(time.Second))
		}
	}()

	syncTimer := func() {
		t, err := client.Sync(ctx, &pb.Profile{
			Name: profile.Name,
		})
		if err != nil {
			log.Panicf("%v", err)
		}

		duration := time.Duration(t.Nanoseconds)
		switch t.GetState() {
		case pb.State_BREAK:
			musicCommand(app, "Pause")
			log.Infof("Break for %s", duration.Truncate(time.Second))
			stateTimer = timer.NewTimer(duration)
		case pb.State_FOCUS:
			musicCommand(app, "Play")
			log.Infof("Focus for %s", duration.Truncate(time.Second))
			stateTimer = timer.NewTimer(duration)
		}
	}

	syncTimer()

	for {
		select {
		case <-stateTimer.C:
			syncTimer()
		case a := <-actions:
			switch a {
			case Play:
				musicCommand(app, "Play")
			case Pause:
				musicCommand(app, "Pause")
			case Toggle:
				musicCommand(app, "PlayPause")
			case Previous:
				musicCommand(app, "Previous")
			case Next:
				musicCommand(app, "Next")
			case Reset:
				musicCommand(app, "Play")
				syncTimer()
			case Remaining:
				d := stateTimer.Remaining()
				log.Infof("Remaining Time: %s", d.Truncate(time.Second).String())
			default:
				// This should never happen
				log.Fatalf("Unrecognized action requested")
			}
		}
	}
}

func StartClient(profile, app, host string, port uint16) {
	err := term.Init()
	if err != nil {
		log.Fatalf("could not start termbox-go: %v", err)
	}
	defer term.Close()

	done := make(chan struct{})
	actions := make(chan Action)

	p := Profile{
		Name: profile,
		Host: host,
		Port: port,
	}
	go Timer(actions, app, p)

	go func() {
		for {
			switch ev := term.PollEvent(); ev.Type {
			case term.EventKey:
				switch ev.Key {
				case term.KeySpace:
					actions <- Toggle
				case term.KeyEsc:
					done <- struct{}{}
				default:
					fmt.Println(ev.Key)
					switch ev.Ch {
					case 'q':
						done <- struct{}{}
					case 'r':
						actions <- Reset
					case 'p':
						actions <- Previous
					case 'n':
						actions <- Next
					default:
						// ignore event
					}
				}
			default:
				// ignore event
			}
		}
	}()

	<-done
	for n := 0; n < 10; n++ {
		time.Sleep(time.Millisecond * 200)
		fmt.Printf("\rShutting down." + strings.Repeat(".", n))
	}
}
