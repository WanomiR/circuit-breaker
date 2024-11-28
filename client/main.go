package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker/v2"
)

const API_URL = "/api/v1/ping"

func main() {
	cb := NewCircuitBreaker("server circuit breaker", 5*time.Second, 3)
	handler := NewHandler(cb)

	http.HandleFunc(API_URL, handler)
	fmt.Println("client app is running")
	http.ListenAndServe(":8080", nil)
}

func NewCircuitBreaker(name string, timeout time.Duration, maxFailures uint32) *gobreaker.CircuitBreaker[[]byte] {
	st := gobreaker.Settings{
		Name: name, Timeout: timeout, ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.TotalFailures >= maxFailures
		},
	}
	return gobreaker.NewCircuitBreaker[[]byte](st)
}

func NewHandler(cb *gobreaker.CircuitBreaker[[]byte]) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := pingToServerApp(cb)

		data := map[string]string{
			"message": "Message",
			"data":    string(res),
		}

		byteData, _ := json.Marshal(data)
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		w.Write(byteData)
	}
}

func pingToServerApp(cb *gobreaker.CircuitBreaker[[]byte]) ([]byte, error) {
	body, err := cb.Execute(doRequets)
	if err != nil {
		return []byte("pong from client"), err
	}

	return body, nil
}

func doRequets() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8082"+API_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error when sending request to the server")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
