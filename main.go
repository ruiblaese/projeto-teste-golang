package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ob "github.com/funkygao/golib/observer"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
)

var (
	listRequests = []string{}
	ctx          = context.Background()
)

func main() {

	testRedis()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"DataFields": listRequests,
		})

	})

	router.GET("/request", func(c *gin.Context) {

		uuid := uuid.Must(uuid.NewV4())
		requestID := fmt.Sprint(uuid)
		log.Print("/request:", requestID)
		listRequests = append(listRequests, requestID)

		eventCh1 := make(chan interface{})
		ob.Subscribe(requestID, eventCh1)
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
					fmt.Println("count: ", len(listRequests), "-> requestId: ", requestID, " -> +1 seg")
					if indexOf(requestID, listRequests) < 0 {
						continueFor = false
					}
				}
			}
		}(c)

	})

	router.GET("/release-request", func(c *gin.Context) {

		requestID := c.DefaultQuery("id", "1")
		indice := indexOf(requestID, listRequests)
		if indice >= 0 {

			ob.Publish(requestID, requestID)
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

func testRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	os.Exit(0)
}
