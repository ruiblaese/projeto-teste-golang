package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	ob "github.com/funkygao/golib/observer"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

var listRequests = []string{}

//teste commit in develop
func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"DataFields": listRequests,
		})

	})

	router.GET("/request", func(c *gin.Context) {

		uuid := uuid.Must(uuid.NewV4())
		requestId := fmt.Sprint(uuid)
		log.Print("/request:", requestId)
		listRequests = append(listRequests, requestId)

		eventCh1 := make(chan interface{})
		ob.Subscribe(requestId, eventCh1)
		func(c *gin.Context) {
			continueFor := true
			for continueFor {
				select {
				case data := <-eventCh1:
					dataS := fmt.Sprint(data)
					c.String(http.StatusOK, dataS)
					continueFor = false
				default:
					time.Sleep(1000 * time.Millisecond)
					c.Writer.Flush()
					fmt.Println("count: ", len(listRequests), "-> requestId: ", requestId, " -> +1 seg")
					if indexOf(requestId, listRequests) < 0 {
						continueFor = false
					}
				}
			}
		}(c)

	})

	router.GET("/release-request", func(c *gin.Context) {

		requestId := c.DefaultQuery("id", "1")
		indice := indexOf(requestId, listRequests)
		if indice >= 0 {

			ob.Publish(requestId, requestId)
			listRequests = append(listRequests[:indice], listRequests[indice+1:]...)
			c.String(http.StatusOK, "OK")

		} else {
			c.String(http.StatusOK, "Error")
		}

	})

	router.Run(":8080")

}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
