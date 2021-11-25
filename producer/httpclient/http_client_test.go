package httpclient

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/rick-and-morty-character-migration/producer/model"
	"github.com/rick-and-morty-character-migration/producer/util"
	"github.com/stretchr/testify/assert"
)

const (
	anyUrl = "www.anyurl.any"

	okResponseBody = `{"info":{"count":826,"pages":42,"next":"https:\/\/rickandmortyapi.com\/api\/character?page=2","prev":null},"results":[{"id":1,"name":"Rick Sanchez","status":"Alive","species":"Human","type":"","gender":"Male","origin":{"name":"Earth (C-137)","url":"https:\/\/rickandmortyapi.com\/api\/location\/1"},"location":{"name":"Citadel of Ricks","url":"https:\/\/rickandmortyapi.com\/api\/location\/3"},"image":"https:\/\/rickandmortyapi.com\/api\/character\/avatar\/1.jpeg","episode":["https:\/\/rickandmortyapi.com\/api\/episode\/1"],"url":"https:\/\/rickandmortyapi.com\/api\/character\/1","created":"2017-11-04T18:48:46.250Z"}]}`
	notSuccessful  = "Not a successful request"

	oneElement = 1
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestGetRequestOk(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), anyUrl)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       util.StrToReadCloser(okResponseBody),
			Header:     make(http.Header),
		}
	})

	httpClientHandler.Client = client

	response, err := httpClientHandler.Get(anyUrl)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var actual model.CharacterResponse
	json.Unmarshal(response.ResponseBody, &actual)
	assert.NotEmpty(t, actual.Results)
	assert.Equal(t, oneElement, len(actual.Results))
}

func TestGetRequestNotFound(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), anyUrl)
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       util.StrToReadCloser(notSuccessful),
			Header:     make(http.Header),
		}
	})

	httpClientHandler.Client = client

	response, err := httpClientHandler.Get(anyUrl)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}
