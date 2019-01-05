package pomodoro

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	pb "github.com/linuxfreak003/go-pomodoro/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Action int

const (
	Start Action = iota
	Stop  Action = iota
	Reset Action = iota
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
	stateTimer := &time.Timer{
		C: make(chan time.Time),
	}

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
			log.Infof("Break for %.2f minutes", duration.Minutes())
			stateTimer = time.NewTimer(duration)
		case pb.State_FOCUS:
			musicCommand(app, "Play")
			log.Infof("Focus for %.2f minutes", duration.Minutes())
			stateTimer = time.NewTimer(duration)
		}
	}

	syncTimer()

	for {
		select {
		case <-stateTimer.C:
			syncTimer()
		case a := <-actions:
			switch a {
			case Start:
				musicCommand(app, "Play")
			case Stop:
				musicCommand(app, "Pause")
			case Reset:
				musicCommand(app, "Play")
				syncTimer()
			}
		}
	}
}

func StartClient(profile, app, host string, port uint16) {
	done := make(chan struct{})
	actions := make(chan Action)

	p := Profile{
		Name: profile,
		Host: host,
		Port: port,
	}
	go Timer(actions, app, p)

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
