package request

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sony/gobreaker/v2"
)

type HttpRequest struct {
	cb     *gobreaker.CircuitBreaker[[]byte]
	client http.Client
	url    string
}

func NewHttpRequest(cb *gobreaker.CircuitBreaker[[]byte], client http.Client, url string) *HttpRequest {
	return &HttpRequest{cb: cb, client: client, url: url}
}

func (h *HttpRequest) PingToServeApp() ([]byte, error) {
	doRequest := func() ([]byte, error) {
		req, err := http.NewRequest(http.MethodGet, h.url, nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := h.client.Do(req)
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

	body, err := h.cb.Execute(doRequest)

	if err != nil {
		return []byte("pong from client"), err
	}

	return body, nil
}
