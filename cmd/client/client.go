package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pedro-makiolke/go-gprc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserBidirectionalStream(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Pedro",
		Email: "P@P.com",
	}
	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Pedro",
		Email: "P@P.com",
	}
	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive: %v", err)
		}
		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "p1",
			Name:  "Pedro",
			Email: "p1@p1.com",
		},
		&pb.User{
			Id:    "p2",
			Name:  "Pedro",
			Email: "p2@p2.com",
		},
		&pb.User{
			Id:    "p3",
			Name:  "Pedro",
			Email: "p3@p3.com",
		},
		&pb.User{
			Id:    "p4",
			Name:  "Pedro",
			Email: "p4@p4.com",
		},
		&pb.User{
			Id:    "p5",
			Name:  "Pedro",
			Email: "p5@p5.com",
		},
		&pb.User{
			Id:    "p6",
			Name:  "Pedro",
			Email: "p6@p6.com",
		},
		&pb.User{
			Id:    "p7",
			Name:  "Pedro",
			Email: "p7@p7.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request:  %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserBidirectionalStream(client pb.UserServiceClient) {

	stream, err := client.AddUserBidirectionalStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "p1",
			Name:  "Pedro",
			Email: "p1@p1.com",
		},
		&pb.User{
			Id:    "p2",
			Name:  "Pedro",
			Email: "p2@p2.com",
		},
		&pb.User{
			Id:    "p3",
			Name:  "Pedro",
			Email: "p3@p3.com",
		},
		&pb.User{
			Id:    "p4",
			Name:  "Pedro",
			Email: "p4@p4.com",
		},
		&pb.User{
			Id:    "p5",
			Name:  "Pedro",
			Email: "p5@p5.com",
		},
		&pb.User{
			Id:    "p6",
			Name:  "Pedro",
			Email: "p6@p6.com",
		},
		&pb.User{
			Id:    "p7",
			Name:  "Pedro",
			Email: "p7@p7.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.GetId())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status: %v", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
