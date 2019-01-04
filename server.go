package pomodoro

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/linuxfreak003/go-pomodoro/pb"
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
	start := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.UTC)
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
				Nanoseconds: (focusPeriod - remaining).Nanoseconds(),
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

func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPomodoroServer(grpcServer, NewServer())
	grpcServer.Serve(lis)
}
