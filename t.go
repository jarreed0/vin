package main

import (
  "context"
  "log"
  "os"

  "google.golang.org/grpc"

  pb "vin.proto"
)

func main() {
  // Set up a connection to the server.
  conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()

  // Create a new client and send a request to the server.
  client := pb.NewMyServiceClient(conn)
  response, err := client.DoSomething(context.Background(), &pb.MyRequest{
    RequestParameter: os.Args[1],
  })
  if err != nil {
    log.Fatalf("could not do something: %v", err)
  }

  // Print the response from the server.
  log.Printf("response: %s", response.ResponseParameter)
}
