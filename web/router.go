package web

import (
	"net/http"
	"time"

	pb "github.com/craigdfrench/event/daemon/grpc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	gopathWebUI = "github.com/craigdfrench/event/web-ui/"
	webUIBuild  = "build/"
)

func WebEngine(static, destination string, db *pb.EventServiceClient) (engine *gin.Engine) {
	engine = gin.Default()
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	// Serve up the React Web GUI
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, static)
	})
	engine.Static("/html", destination)
	engine.GET("/event", FormReadEvents(*db))
	engine.POST("/event", FormWriteEvent(*db))
	return
}
