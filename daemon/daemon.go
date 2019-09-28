package daemon 

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "github.com/craigdfrench/event/daemon/grpc"
	"github.com/craigdfrench/event/storage"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
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
	id, err := insertEvent(s.Database, in.CreatedAt, in.Email, in.Environment, in.Component, in.Message, in.Data)
	return &pb.EventIdentifier{Id: id}, err
}

// WriteEvent implements storage.QueryMultipleEvents
func (s *EventServer) QueryMultipleEvents(ctx context.Context, in *pb.QueryEventRequest) (response *pb.QueryEventResponse, err error) {
	//log.Println("Received: ", in.CreatedAt, in.Email, in.Environment, in.Component, in.Message, in.Data)
	startTime := time.Now()
	if startTime, err = ptypes.Timestamp(in.GetTimeRange().GetStartTime()); err != nil {
		startTime = time.Now()
	}
	eventList := []*pb.Event{}
	query := Event{
		CreatedAt:   startTime,
		Email:       in.Email,
		Environment: in.Environment,
		Component:   in.Component,
		Message:     in.Message}
	eventList, err = getEvents(s.Database, query)
	response = &pb.QueryEventResponse{Results: eventList}
	return
}

// WriteEvent implements storage.QueryMultipleEvents
func (s *EventServer) ReadSingleEvent(ctx context.Context, in *pb.EventIdentifier) (*pb.Event, error) {
	//log.Println("Received: ", in.CreatedAt, in.Email, in.Environment, in.Component, in.Message, in.Data)
	event := pb.Event{}
	return &event, nil
}

