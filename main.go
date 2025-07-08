package main

import (
	"logview_ui/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := routes.SetupRouter()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/logview/pages/logtable")
	})
	// Listen and Server in 0.0.0.0:8057
	r.Run(":8057")
}
