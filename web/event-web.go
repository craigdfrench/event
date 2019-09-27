package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/craigdfrench/event-service/contracts"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	gopathWebUI = "github.com/craigdfrench/event-logger/web-ui/"
	webUIBuild  = "build/"
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

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	db := pb.NewStorageServiceClient(conn)

	var htmlPath string

	if golangpath, present := os.LookupEnv("GOPATH"); present {
		htmlPath = golangpath + "/src/" + gopathWebUI + webUIBuild
	} else {
		htmlPath = "./" + webUIBuild
	}

	r := gin.Default()

	// Serve up the React Web GUI
	r.GET("/", func(c *gin.Context) {
		log.Printf("Web directory is %s", htmlPath)

		c.Redirect(http.StatusTemporaryRedirect, "html/index.html")
	})
	r.Static("/html", htmlPath)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": htmlPath,
		})
	})
	r.GET("/event", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		request := &pb.QueryEventRequest{
			Email:       c.Query("Email"),
			Environment: c.Query("Environment"),
			Component:   c.Query("Component"),
			Message:     c.Query("Message")}

		eventRecords, err := db.QueryMultipleEvents(ctx, request)
		if err != nil {
			panic(err.Error())
		}
		c.JSON(http.StatusOK, eventRecords.GetResults())
	})
	r.POST("/event", func(c *gin.Context) {
		var form Event
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(form)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := db.WriteSingleEvent(ctx, &pb.Event{
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
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
