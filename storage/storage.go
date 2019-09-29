package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	pb "github.com/craigdfrench/event/daemon/grpc"

	// Tied to postgreSQL
	_ "github.com/lib/pq"
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

// Server is used to implement storage.EventServiceServer
type Server struct {
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

// InsertEvent will insert the record format
func InsertEvent(db *sql.DB, CreatedAt, Email, Environment, Component, Message, Data string) (ID string, err error) {
	sqlStatement := `
		INSERT INTO public.event ("CreatedAt", "Email", "Environment", "Component", "Message", "Data")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING "Id" `
	ID = ""
	err = db.QueryRow(sqlStatement, CreatedAt, Email, Environment, Component, Message, Data).Scan(&ID)
	return
}

// GetEvents will retrieve events as per query
func GetEvents(db *sql.DB, query Event) ([]*pb.Event, error) {
	fmt.Printf("query is %s %d %d", query, len(query.Message), len(query.Environment))
	var queryString []string
	var queryArgs []interface{}
	if len(query.Component) > 0 {
		queryString = append(queryString, fmt.Sprintf(`"Component" = $%d`, len(queryArgs)+1))
		queryArgs = append(queryArgs, query.Component)
	}
	if len(query.Email) > 0 {
		queryString = append(queryString, fmt.Sprintf(`"Email" = $%d`, len(queryArgs)+1))
		queryArgs = append(queryArgs, query.Email)
	}
	if len(query.Environment) > 0 {
		queryString = append(queryString, fmt.Sprintf(`"Environment" = $%d`, len(queryArgs)+1))
		queryArgs = append(queryArgs, query.Environment)
	}
	if len(query.Message) > 0 {
		queryString = append(queryString, fmt.Sprintf(`POSITION($%d in "Message")>0`, len(queryArgs)+1))
		queryArgs = append(queryArgs, query.Message)
	}
	fmt.Println("queryString is ", queryString, "=>", queryArgs)
	var whereClause string
	switch len(queryString) {
	case 0:
		whereClause = "TRUE"
	case 1:
		whereClause = queryString[0]
	default:
		whereClause = strings.Join(queryString, " AND ")
	}
	fmt.Printf("whereClause  is SELECT * from public.event WHERE %s", whereClause)
	rows, err := db.Query("SELECT * from public.event WHERE "+whereClause, queryArgs...)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("In getEvents", rows)
	eventRecords := []*pb.Event{}
	for rows.Next() {
		eventRecord := pb.Event{}
		if err = rows.Scan(&eventRecord.Id, &eventRecord.CreatedAt, &eventRecord.Email, &eventRecord.Environment, &eventRecord.Component, &eventRecord.Message, &eventRecord.Data); err != nil {
			fmt.Println("Errored out", err.Error())

			eventRecords = nil
			break
		}
		fmt.Println(eventRecord)
		eventRecords = append(eventRecords, &eventRecord)
	}
	return eventRecords, err
}

