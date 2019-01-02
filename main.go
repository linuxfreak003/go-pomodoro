package pomodoro

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
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

var layout = "06-02-01 15:04:05"

func startMusic(app string) {
	log.Infof("[%v] Starting music on %s", time.Now().Format(layout), app)
	cmdString := fmt.Sprintf("tell app \"%s\" to play", app)
	cmd := exec.Command("osascript", "-e", cmdString)
	err := cmd.Run()
	if err != nil {
		log.Errorf("%v", err)
	}
}

func stopMusic(app string) {
	log.Infof("[%v] Stopping music on %s", time.Now().Format(layout), app)
	cmdString := fmt.Sprintf("tell app \"%s\" to pause", app)
	cmd := exec.Command("osascript", "-e", cmdString)
	err := cmd.Run()
	if err != nil {
		log.Errorf("%v", err)
	}
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
			stopMusic(app)
			log.Infof("Break for %d minutes", int(t.Duration))
			breakTimer = time.NewTimer(time.Duration(t.Duration) * time.Minute)
		}
		if t.GetState() == pb.State_FOCUS {
			startMusic(app)
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
				startMusic(app)
			case Stop:
				stopMusic(app)
			case Reset:
				startMusic(app)
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
