package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	ob "github.com/funkygao/golib/observer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var listRequests = []string{}

func main() {

	// Echo instance
	router := echo.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	//router.LoadHTMLGlob("templates/*")

	router.GET("/", index)

	router.GET("/request", request)

	router.GET("/release-request", releaseRequest)

	// Start server
	router.Logger.Fatal(router.Start(":8080"))

}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func request(c echo.Context) error {

	//uuid := uuid.Must(uuid.NewV4())
	requestId := RandStringBytes(32)
	log.Print("/request:", requestId)
	listRequests = append(listRequests, requestId)

	eventCh1 := make(chan interface{})
	ob.Subscribe(requestId, eventCh1)
	return func(c echo.Context) error {

		for {
			select {
			case <-eventCh1:
				fmt.Println("tick.")
				data := <-eventCh1
				dataS := fmt.Sprint(data)
				return c.String(http.StatusOK, dataS)
			default:
				time.Sleep(1000 * time.Millisecond)
				c.Response().Flush()
				fmt.Println("count: ", len(listRequests), "-> requestId: ", requestId, " -> +1 seg")
			}
		}

	}(c)

}

func releaseRequest(c echo.Context) error {

	requestId := c.QueryParam("id")
	indice := indexOf(requestId, listRequests)
	if indice >= 0 {

		ob.Publish(requestId, requestId)
		listRequests = append(listRequests[:indice], listRequests[indice+1:]...)
		return c.String(http.StatusOK, "OK")

	} else {
		return c.String(http.StatusOK, "Error")
	}

}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
