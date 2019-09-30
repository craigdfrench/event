package web

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/craigdfrench/event/daemon/grpc"
	"github.com/gin-gonic/gin"
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

		request := &pb.QueryEventRequest{
			Email:       c.Query("Email"),
			Environment: c.Query("Environment"),
			Component:   c.Query("Component"),
			Message:     c.Query("Message")}

		eventRecords, err := service.QueryMultipleEvents(ctx, request)
		if err != nil {
			panic(err.Error())
		}
		c.JSON(http.StatusOK, eventRecords.GetResults())
	}
}
