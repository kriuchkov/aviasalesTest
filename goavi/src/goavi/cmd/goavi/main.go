package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"goavi/pkg/storage"
)

func main() {
	srv := &server{
		storage: storage.NewStorage(),
		http:    &http.Server{Addr: ":8080"},
	}

	srv.http.Handler = srv.setupRoute()
	srv.http.ListenAndServe()
}

type server struct {
	storage *storage.Storage
	http    *http.Server
}

func (srv *server) setupRoute() *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.POST("/receive", srv.receiveXML)
	router.GET("/itinerary", srv.getFlight)
	return router
}

func (srv *server) receiveXML(c *gin.Context) {
	var resp Response
	var err error

	if c.Request.Body == nil {
		resp.Error = "Запрос пуст"
		c.IndentedJSON(401, resp)
		c.Abort()
	}

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body.Close()

	err = srv.storage.LoadXML(bodyBytes)

	if err != nil {
		resp.Status = "error"
		resp.Error = fmt.Sprintf("Error %v", err)
	} else {
		resp.Status = "ok"
	}

	payload, _ := json.Marshal(resp)
	c.Data(http.StatusOK, "text/json", payload)
}

type Query struct {
	Source      string `form:"source"`
	Destination string `form:"destination"`

	// Необязательные параметры
	Type   string `form:"type"`
	Return bool   `form:"return"`
}

func (srv *server) getFlight(c *gin.Context) {
	var resp Response

	query := new(Query)
	if err := c.BindQuery(query); err != nil {
		resp.Error = fmt.Sprintf("%v", err)
		c.IndentedJSON(401, resp)
		c.Abort()
	}

	out := srv.storage.GetItinerary(query.Source, query.Destination, query.Return)
	if len(out) == 0 {
		resp.Error = "Запрос пуст"
		c.IndentedJSON(401, resp)
		c.Abort()
	}

	if query.Type != "" {
		var q storage.StorageList
		switch query.Type {
		case "stime": //самый долгий маршрут
			q = storage.NewTimeQueueMax()
		case "ltime": // самый быстрый маршрут
			q = storage.NewTimeQueueMin()
		case "sprice": // самый дорогой маршрут
			q = storage.NewPriceQueueMax()
		case "lprice": // самый дешевый маршрут
			q = storage.NewPriceQueueMin()
		case "optimal": // оптимальный маршрут
			q = storage.NewOptimalQueue()
		}
		srv.storage.OptimalItinerary(out, q)
		resp.Result = q.PopOrdered()
	}

	resp.Status = "ok"
	payload, _ := json.Marshal(resp)
	c.Data(http.StatusOK, "text/json", payload)
}

type Response struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}
