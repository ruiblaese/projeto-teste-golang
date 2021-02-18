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

var ctx = context.Background()

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"DataFields": getAllRequests(rdb),
		})

	})

	router.GET("/request", func(c *gin.Context) {

		uuid := uuid.Must(uuid.NewV4())
		requestID := fmt.Sprint(uuid)
		log.Print("/request:", requestID)
		addRequest(rdb, requestID)

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
					fmt.Println("count: ", len(getAllRequests(rdb)), "-> requestId: ", requestID, " -> +1 seg")
					if !getRequest(rdb, requestID) {
						continueFor = false
					}
				}
			}
		}(c)

	})

	router.GET("/release-request", func(c *gin.Context) {

		requestID := c.DefaultQuery("id", "1")
		data := c.DefaultQuery("data", "")
		if getRequest(rdb, requestID) {

			ob.Publish(requestID, data)
			deleteRequest(rdb, requestID)
			c.String(http.StatusOK, "OK")

		} else {
			c.String(http.StatusOK, "Error")
		}

	})

	router.Run(":8080")

}

func addRequest(rdb *redis.Client, requestID string) {

	err := rdb.Set(ctx, requestID, "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

}

func getRequest(rdb *redis.Client, requestID string) bool {
	_, err := rdb.Get(ctx, requestID).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		panic(err)
	} else {
		return true
	}
}

func getAllRequests(rdb *redis.Client) []string {
	keys := rdb.Keys(ctx, "*")
	s, _ := keys.Result()
	return s
}

func deleteRequest(rdb *redis.Client, requestID string) {
	rdb.Del(ctx, requestID)

}
