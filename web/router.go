package web

import (
	"net/http"

	pb "github.com/craigdfrench/event/daemon/grpc"
	"github.com/gin-gonic/gin"
)

const (
	gopathWebUI = "github.com/craigdfrench/event/web-ui/"
	webUIBuild  = "build/"
)

func WebEngine(static, destination string, db *pb.EventServiceClient) (engine *gin.Engine) {
	engine = gin.Default()

	// Serve up the React Web GUI
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, static)
	})
	engine.Static("/html", destination)
	engine.GET("/event", FormReadEvents(*db))
	engine.POST("/event", FormWriteEvent(*db))
	return
}
