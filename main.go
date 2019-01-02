package pomodoro

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
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

func Timer(actions chan Action, app, profile string) {
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		log.Errorf("%v", err)
	}
	defer conn.Close()
	client := pb.NewPomodoroClient(conn)
	focusTimer := &time.Timer{
		C: make(chan time.Time),
	}
	breakTimer := &time.Timer{
		C: make(chan time.Time),
	}

	syncTimer := func() {
		t, err := client.Sync(ctx, &pb.Profile{
			Name: profile,
		})
		if err != nil {
			log.Errorf("%v", err)
		}

		if t.GetState() == pb.State_BREAK {
			musicCommand(app, "Pause")
			log.Infof("Break for %d minutes", int(t.Duration))
			breakTimer = time.NewTimer(time.Duration(t.Duration) * time.Minute)
		}
		if t.GetState() == pb.State_FOCUS {
			musicCommand(app, "Play")
			log.Infof("Focus for %d minutes", int(t.Duration))
			focusTimer = time.NewTimer(time.Duration(t.Duration) * time.Minute)
		}
	}

	syncTimer()

	for {
		select {
		case <-focusTimer.C:
			syncTimer()
		case <-breakTimer.C:
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

func StartClient() {
	var profile, app string

	flag.StringVar(&app, "app", "spotify", "music app to use")
	flag.StringVar(&profile, "profile", "Default", "profile to sync with")
	flag.Parse()

	done := make(chan struct{})
	actions := make(chan Action)

	go Timer(actions, app, profile)

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

type Server struct{}

func DefaultProfileTime() *pb.Timer {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.UTC)
	remaining := now.Sub(start).Minutes()
	focusPeriod := float64(25)
	breaks := []float64{5, 5, 5, 15}
	breakIndex := 0
	for {
		if remaining < focusPeriod {
			return &pb.Timer{
				Duration: focusPeriod - remaining,
				State:    pb.State_FOCUS,
			}
		}

		remaining -= focusPeriod

		curBreak := breaks[breakIndex]
		if remaining < curBreak {
			return &pb.Timer{
				Duration: focusPeriod - remaining,
				State:    pb.State_BREAK,
			}
		}

		remaining -= curBreak

		breakIndex++
		if breakIndex >= len(breaks) {
			breakIndex = 0
		}
	}
}

func (s *Server) Sync(ctx context.Context, req *pb.Profile) (*pb.Timer, error) {
	// Default Profile
	return DefaultProfileTime(), nil
}

func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPomodoroServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}
