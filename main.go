package main

import (
	"fmt"
	"net/http"

	ob "github.com/funkygao/golib/observer"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/request", func(c *gin.Context) {

		requestId := c.DefaultQuery("id", "1")

		eventCh1 := make(chan interface{})
		ob.Subscribe(requestId, eventCh1)
		func(c *gin.Context) {
			data := <-eventCh1
			dataS := fmt.Sprint(data)
			c.String(http.StatusOK, dataS)
		}(c)

	})

	router.GET("/release-request", func(c *gin.Context) {

		requestId := c.DefaultQuery("id", "1")

		ob.Publish(requestId, requestId)
		c.String(http.StatusOK, "OK")
	})

	router.Run(":8080")

}
