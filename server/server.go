package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/bluele/slack"
	pb "github.com/linuxfreak003/go-pomodoro/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	Profiles map[*pb.Profile]*pb.Timer
}

func NewServer() *Server {
	return &Server{
		Profiles: make(map[*pb.Profile]*pb.Timer),
	}
}

func DefaultProfileTime() *pb.Timer {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.UTC)
	remaining := now.Sub(start) // Just leave as a time.Duration

	Minutes := func(m int) time.Duration {
		return time.Duration(m) * time.Minute
	}

	focusPeriod := Minutes(25)
	breaks := []time.Duration{Minutes(5), Minutes(5), Minutes(5), Minutes(15)}
	breakIndex := 0
	for {
		if remaining < focusPeriod {
			return &pb.Timer{
				Nanoseconds: (focusPeriod - remaining).Nanoseconds(),
				State:       pb.State_FOCUS,
			}
		}

		remaining -= focusPeriod

		curBreak := breaks[breakIndex]
		if remaining < curBreak {
			return &pb.Timer{
				Nanoseconds: (curBreak - remaining).Nanoseconds(),
				State:       pb.State_BREAK,
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
	if timer, ok := s.Profiles[req]; ok {
		return timer, nil
	}
	return DefaultProfileTime(), nil
}

func StartServer(port uint16, token, channelName string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	host := GetOutboundIP()
	addrMsg := fmt.Sprintf("Server running, connect with `go-pomodoro client --host %s --port %d`", host, port)
	log.Infof(addrMsg)

	api := slack.New(token)
	err = api.ChatPostMessage(channelName, addrMsg, &slack.ChatPostMessageOpt{
		AsUser: true,
	})
	if err != nil {
		log.Warnf("slack message not sent: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPomodoroServer(grpcServer, NewServer())
	grpcServer.Serve(lis)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
