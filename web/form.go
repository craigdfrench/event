package web

import (
	"context"
	"log"
	"net/http"
	"time"
	pb "github.com/craigdfrench/event/daemon/grpc"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// Event generated from form
type Event struct {
	ID          string    `form:"id"`
	Email       string    `form:"email" binding:"required"`
	CreatedAt   time.Time `form:"createdAt"`
	Environment string    `form:"environment" binding:"required"`
	Component   string    `form:"component" binding:"required"`
	Message     string    `form:"message"`
	Data        string    `form:"data"`
}

func FormWriteEvent(service pb.EventServiceClient) func(c *gin.Context) {

	return func(c *gin.Context) {
		var form Event
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := service.WriteSingleEvent(ctx, &pb.Event{
			CreatedAt:   form.CreatedAt.String(),
			Email:       form.Email,
			Environment: form.Environment,
			Component:   form.Component,
			Message:     form.Message,
			Data:        form.Data})
		if err != nil {
			log.Fatalf("could not write: %v", err)
		}
		log.Printf("Generated ID: %s", r.GetId())
		c.JSON(http.StatusOK, gin.H{"id": r.GetId()})
	}
}

func FormReadEvents(service pb.EventServiceClient) func(*gin.Context) {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var startTime,endTime time.Time
		var startTimeStamp, endTimeStamp *timestamp.Timestamp
		var err error

		if startTime, err = time.Parse("2006-01-02T15:04:05-07:00", c.Query("startTime")); err == nil {
			if startTimeStamp, err = ptypes.TimestampProto(startTime); err != nil {
				startTimeStamp = nil
				log.Printf("Unable to parse startTime: %s %s",c.Query("startTime"),err.Error())
			} 
		} else {
			log.Printf("Unable to parse startTime: %s %s",c.Query("startTime"),err.Error())
		}
		if endTime, err = time.Parse("2006-01-02T15:04:05-07:00", c.Query("endTime")); err == nil {
			if endTimeStamp, err = ptypes.TimestampProto(endTime); err != nil {
				endTimeStamp = nil
				log.Printf("Unable to parse endTime: %s %s",c.Query("endTime"),err.Error())
			}
		} else {
			log.Printf("Unable to parse endTime: %s %s",c.Query("endTime"),err.Error())
		}
		request := &pb.QueryEventRequest{
			Email:       c.Query("email"),
			Environment: c.Query("environment"),
			Component:   c.Query("component"),
			Message:     c.Query("message"),
			TimeRange: &pb.TimeQuery{
				StartTime: startTimeStamp,
				EndTime:   endTimeStamp,
				Duration:  nil,
			},
		}
		eventRecords, err := service.QueryMultipleEvents(ctx, request)
		if err != nil {
			panic(err.Error())
		}
		c.JSON(http.StatusOK, eventRecords.GetResults())
	}
}
