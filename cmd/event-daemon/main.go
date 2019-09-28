package main

import (
	"log"
	"net"
	"os"

	"github.com/craigdfrench/event/daemon"
	pb "github.com/craigdfrench/event/daemon/grpc"
	"github.com/craigdfrench/event/storage"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
	// EventDatabaseBackend type
	EventDatabaseBackend = "postgres"
	// EventTableDefinition for definition of event table
	EventTableDefinition = "event.relation.sql"
	// GoPathSrcDir is where the event definition is found
	GoPathSrcDir = "/src/github.com/craigdfrench/event-service/storage/"
	// EventDatabaseConnectionString specifies credentials to access database
	EventDatabaseConnectionString = "user=pqgotest dbname=pqgotest password=pqgotest sslmode=disable"
)

func main() {
	var schema string
	if gopath, present := os.LookupEnv("GOPATH"); present {
		schema = gopath + GoPathSrcDir + EventTableDefinition
	} else {
		schema = "./" + EventTableDefinition
	}

	db, err := storage.SetupDatabase(EventDatabaseBackend, EventDatabaseConnectionString, schema)
	if err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEventServiceServer(s, &daemon.EventServer{db})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
