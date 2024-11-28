package main

import (
	"client/config"
	"client/handler"
	"client/request"
	"fmt"
	"net/http"
)

func main() {
	cb := config.CircuitBreakerConfig()
	req := request.NewHttpRequest(cb, http.Client{}, "http://localhost:8082/api/v1/ping")
	hdl := handler.NewHandler(req)

	http.HandleFunc("/api/v1/ping", hdl.Ping)
	fmt.Println("client app is running")
	http.ListenAndServe(":8080", nil)
}
