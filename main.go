package main

import (
	"fmt"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/alertmanager/types"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"time"

)

type status string
type errorType string

const (
	statusCodeSuccess status = "success"
	statusCodeError   status = "error"
		duration      = 30 * time.Second
		interval_time = 5 * time.Second
)

type response struct {
	Status    status      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType errorType   `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func main() {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./views.html", true)))
	ticker := time.NewTicker(interval_time)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("repeat task....")

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	api := router.Group("/api")
	{
		api.POST("/v1/alerts", postAlerts)
	}

	router.Run(":3909")
	
}

func postAlerts(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var alerts []*types.Alert
	var vaildAlerts []*types.Alert
	c.BindJSON(&alerts)
	//responseSucces(c, alerts)
	for _, alert := range alerts{
		fmt.Println(alert.Labels)
		if alert.Labels["resolved_by"] != "" {
			fmt.Println(alert.Labels["resolved_by"])
			continue
			/* if isResolved(alert){
				checkResolved(alert)
				continue
			} */
		}
		vaildAlerts = append(vaildAlerts, alert)

	}
	forwardAlerts(vaildAlerts)
}
func forwardAlerts(alerts []*types.Alert){
	b, err := json.Marshal(alerts)
	if err != nil {
		panic(err)
	}
	url := "http://192.168.60.100:9093/api/v1/alerts"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("X-Custom-Header", "resolvedAlerts")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

	}


func responseSucces(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application")
	c.JSON(200, response{
		Status: statusCodeSuccess,
		Data:   data,
	})
}

func isResolved(alert *types.Alert) bool{
	if alert.EndsAt.After(time.Now()){
		return false
	} 
	return true
}
/* func checkResolved(alert *type.Alert){
	
} */
//
// func ()  {
//
// }
