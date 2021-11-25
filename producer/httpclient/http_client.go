package httpclient

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rick-and-morty-character-migration/producer/model"
)

const (
	ContentTypeHeader = "Content-Type"

	ContentTypeJSON = "application/json"
)

var httpClientHandler HttpClientHandler

type HttpClient interface {
	Get(url string) (model.HttpClientResponse, error)
}

type HttpClientHandler struct {
	Client *http.Client
}

func NewHttpHandler(client *http.Client) HttpClientHandler {
	if httpClientHandler.Client == nil {
		httpClientHandler = HttpClientHandler{
			Client: client,
		}
	}

	return httpClientHandler
}

func (httpClientHandler HttpClientHandler) Get(url string) (model.HttpClientResponse, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return model.HttpClientResponse{}, err
	}

	httpClientResponse, err := httpClientHandler.do(request)
	if err != nil {
		return model.HttpClientResponse{}, err
	}

	return httpClientResponse, nil
}

func (httpClientHandler HttpClientHandler) do(request *http.Request) (model.HttpClientResponse, error) {
	request.Header.Add(ContentTypeHeader, ContentTypeJSON)

	log.Println("Sending request to ", request.URL.String())
	response, err := httpClientHandler.Client.Do(request)
	if err != nil {
		return model.HttpClientResponse{}, err
	}
	defer response.Body.Close()

	responseBodyParsed, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model.HttpClientResponse{}, nil
	}

	httpResponse := model.HttpClientResponse{
		ResponseBody: responseBodyParsed,
		StatusCode:   response.StatusCode,
	}

	return httpResponse, nil
}
