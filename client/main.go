package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	chat "github.com/mashardi21/HippoChat/proto"
	"google.golang.org/grpc"
)

var client chat.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

// Connect manages the stream that recieves messages from the server
// and displays them when they are recieved.
func Connect(user *chat.User) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &chat.Connect{User: user, Active: true})

	if err != nil {
		return fmt.Errorf("Connection failed %v", err)
	}

	wait.Add(1)

	go func(str chat.Broadcast_CreateStreamClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()

			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}

			fmt.Printf("%v: %s\n", msg.User.UserName, msg.Body)
		}
	}(stream)

	return streamerror
}

func main() {
	timestamp := time.Now()
	done := make(chan int)

	userName := flag.String("u", "Egbog", "The username that will be used for the current session")
	id := flag.String("i", "1", "This flag is only temporary and will not be present in future versions")
	flag.Parse()

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect to the server: %v", err)
	}

	client = chat.NewBroadcastClient(conn)
	user := &chat.User{
		UserName: *userName,
		ID:       *id,
	}

	Connect(user)

	wait.Add(1)

	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			msg := &chat.Message{
				User:      user,
				Body:      scanner.Text(),
				TimeStamp: timestamp.String(),
			}

			_, err := client.BroadcastMessage(context.Background(), msg)

			if err != nil {
				fmt.Printf("Error sending message: %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
