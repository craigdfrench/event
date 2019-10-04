package daemon

import (
	"context"
	"database/sql"
	"log"
	"time"

	pb "github.com/craigdfrench/event/daemon/grpc"
	"github.com/craigdfrench/event/storage"
)

const (
	port = ":50051"
	// EventDatabaseBackend type
	EventDatabaseBackend = "postgres"
	// EventTableDefinition for definition of event table
	EventTableDefinition = "event.relation.sql"
	// GoPathSrcDir is where the event definition is found
	GoPathSrcDir = "/src/github.com/craigdfrench/event/storage/"
	// EventDatabaseConnectionString specifies credentials to access database
	EventDatabaseConnectionString = "user=pqgotest dbname=pqgotest password=pqgotest sslmode=disable"
)

// EventServer is used to implement daemon
type EventServer struct {
	Database *sql.DB
}

// Event structure
type Event struct {
	ID          string
	Email       string
	CreatedAt   time.Time
	Environment string
	Component   string
	Message     string
	Data        string
}

// WriteEvent implements storage.WriteEvent
func (s *EventServer) WriteSingleEvent(ctx context.Context, in *pb.Event) (*pb.EventIdentifier, error) {
	log.Println("Received: ", in.CreatedAt, in.Email, in.Environment, in.Component, in.Message, in.Data)
	id, err := storage.InsertEvent(s.Database, in.CreatedAt, in.Email, in.Environment, in.Component, in.Message, in.Data)
	return &pb.EventIdentifier{Id: id}, err
}

// WriteEvent implements storage.QueryMultipleEvents
func (s *EventServer) QueryMultipleEvents(ctx context.Context, in *pb.QueryEventRequest) (response *pb.QueryEventResponse, err error) {
	//log.Println("Received: ", in.CreatedAt, in.Email, in.Environment, in.Component, in.Message, in.Data)
	eventList := []*pb.Event{}
	eventList, err = storage.GetEvents(s.Database, *in)
	response = &pb.QueryEventResponse{Results: eventList}
	return
}

// WriteEvent implements storage.QueryMultipleEvents
func (s *EventServer) ReadSingleEvent(ctx context.Context, in *pb.EventIdentifier) (*pb.Event, error) {
	//log.Println("Received: ", in.CreatedAt, in.Email, in.Environment, in.Component, in.Message, in.Data)
	event := pb.Event{}
	return &event, nil
}
