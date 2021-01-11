package main

import (
	"fmt"
	"log"
	"net/http"

	ob "github.com/funkygao/golib/observer"
	"github.com/gofrs/uuid"
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

	uuid := uuid.Must(uuid.NewV4())
	requestId := fmt.Sprint(uuid)
	log.Print("/request:", requestId)
	listRequests = append(listRequests, requestId)

	eventCh1 := make(chan interface{})
	ob.Subscribe(requestId, eventCh1)
	return func(c echo.Context) error {
		data := <-eventCh1
		dataS := fmt.Sprint(data)
		return c.String(http.StatusOK, dataS)
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
