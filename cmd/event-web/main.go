package main

import (
	"log"
	"os"

	"github.com/craigdfrench/event/web"

	pb "github.com/craigdfrench/event/daemon/grpc"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	gopathWebUI = "github.com/craigdfrench/event/web-ui/"
	webUIBuild  = "build/"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	service := pb.NewEventServiceClient(conn)

	var htmlPath string

	if golangpath, present := os.LookupEnv("GOPATH"); present {
		htmlPath = golangpath + "/src/" + gopathWebUI + webUIBuild
	} else {
		htmlPath = "./" + webUIBuild
	}

	engine := web.WebEngine("/html", htmlPath, &service)
	engine.Run() // listen and serve on 0.0.0.0:8080
}
