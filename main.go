package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	chat "github.com/mashardi21/HippoChat/proto"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

// Connection is used to store the data for each connected client.
type Connection struct {
	stream   chat.Broadcast_CreateStreamServer
	userName string
	active   bool
	error    chan error
}

// Server is used to store a slice of all connected clients
type Server struct {
	Connections []*Connection
	chat.UnimplementedBroadcastServer
}

// CreateStream takes a client connection and correlates it with a stream.
// It then adds this connection to the server's list of all connected clients
func (s *Server) CreateStream(pconn *chat.Connect, stream chat.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream:   stream,
		userName: pconn.User.UserName,
		active:   true,
		error:    make(chan error),
	}

	s.Connections = append(s.Connections, conn)

	return <-conn.error
}

// BroadcastMessage takes a new message and sends it to the stream associated with each
// client currently active on the server.
func (s *Server) BroadcastMessage(ctx context.Context, msg *chat.Message) (*chat.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connections {
		wait.Add(1)

		go func(msg *chat.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				grpcLog.Info("Sending message to: ", conn.userName)

				if err != nil {
					log.Fatalf("Could not send message: %v", err)
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &chat.Close{}, nil
}

func main() {
	var connections []*Connection

	server := &Server{Connections: connections}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Error creating the server: %v", err)
	}

	grpcLog.Info("Starting server on port :8080")

	chat.RegisterBroadcastServer(grpcServer, server)
	grpcServer.Serve(listener)
}
