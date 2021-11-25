package fetcher

import (
	"fmt"
	"net/http"
	"testing"

	_ "github.com/rick-and-morty-character-migration/producer/httpclient"
	"github.com/rick-and-morty-character-migration/producer/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	oneCharacterResponse        = `{"info":{"count":826,"pages":42,"next":"","prev":null},"results":[{"id":1,"name":"Rick Sanchez","status":"Alive","species":"Human","type":"","gender":"Male","origin":{"name":"Earth (C-137)","url":"https:\/\/rickandmortyapi.com\/api\/location\/1"},"location":{"name":"Citadel of Ricks","url":"https:\/\/rickandmortyapi.com\/api\/location\/3"},"image":"https:\/\/rickandmortyapi.com\/api\/character\/avatar\/1.jpeg","episode":["https:\/\/rickandmortyapi.com\/api\/episode\/1"],"url":"https:\/\/rickandmortyapi.com\/api\/character\/1","created":"2017-11-04T18:48:46.250Z"}]}`
	paginationCharacterResponse = `{"info":{"count":826,"pages":42,"next":"/character?page=2","prev":null},"results":[{"id":1,"name":"Rick Sanchez","status":"Alive","species":"Human","type":"","gender":"Male","origin":{"name":"Earth (C-137)","url":"https:\/\/rickandmortyapi.com\/api\/location\/1"},"location":{"name":"Citadel of Ricks","url":"https:\/\/rickandmortyapi.com\/api\/location\/3"},"image":"https:\/\/rickandmortyapi.com\/api\/character\/avatar\/1.jpeg","episode":["https:\/\/rickandmortyapi.com\/api\/episode\/1"],"url":"https:\/\/rickandmortyapi.com\/api\/character\/1","created":"2017-11-04T18:48:46.250Z"}]}`
	oneLocationResponse         = `{"info":{"count":126,"pages":7,"next":null,"prev":null},"results":[{"id":1,"name":"Earth (C-137)","type":"Planet","dimension":"Dimension C-137","residents":["https:\/\/rickandmortyapi.com\/api\/character\/38"],"url":"https:\/\/rickandmortyapi.com\/api\/location\/1","created":"2017-11-10T12:42:04.162Z"}]}`
	paginationLocationResponse  = `{"info":{"count":126,"pages":7,"next":"\/location?page=2","prev":null},"results":[{"id":1,"name":"Earth (C-137)","type":"Planet","dimension":"Dimension C-137","residents":["https:\/\/rickandmortyapi.com\/api\/character\/38"],"url":"https:\/\/rickandmortyapi.com\/api\/location\/1","created":"2017-11-10T12:42:04.162Z"}]}`

	notAJSONResponseBody = "not a json"

	charactersPageOneUrl = "/character?page=1"
	charactersPageTwoUrl = "/character?page=2"

	locationsPageOneUrl = "/location?page=1"
	locationsPageTwoUrl = "/location?page=2"

	oneElement  = 1
	twoElements = 2
)

type httpClientMock struct {
	mock.Mock
}

func (m httpClientMock) Get(url string) (model.HttpClientResponse, error) {
	args := m.Called(url)

	return args.Get(0).(model.HttpClientResponse), args.Error(1)
}

func TestFetchAllCharactersWhenOneElementIsReturned(t *testing.T) {
	response := model.HttpClientResponse{
		ResponseBody: []byte(oneCharacterResponse),
		StatusCode:   http.StatusOK,
	}
	mock := new(httpClientMock)

	mock.On("Get", charactersPageOneUrl).Return(response, nil)

	dataFetcher.httpClientHandler = mock

	actual, err := dataFetcher.FetchAllCharacters()

	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
	assert.Equal(t, oneElement, len(actual))
}

func TestFetchAllCharactersPagination(t *testing.T) {
	responsePageOne := model.HttpClientResponse{
		ResponseBody: []byte(paginationCharacterResponse),
		StatusCode:   http.StatusOK,
	}
	responsePageTwo := model.HttpClientResponse{
		ResponseBody: []byte(oneCharacterResponse),
		StatusCode:   http.StatusOK,
	}
	mock := new(httpClientMock)

	mock.On("Get", charactersPageOneUrl).Once().Return(responsePageOne, nil)
	mock.On("Get", charactersPageTwoUrl).Once().Return(responsePageTwo, nil)

	dataFetcher.httpClientHandler = mock

	actual, err := dataFetcher.FetchAllCharacters()

	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
	assert.Equal(t, twoElements, len(actual))
}

func TestFetchAllCharactersWhenUnmarshalError(t *testing.T) {
	response := model.HttpClientResponse{
		ResponseBody: []byte(notAJSONResponseBody),
		StatusCode:   http.StatusOK,
	}
	mock := new(httpClientMock)

	mock.On("Get", charactersPageOneUrl).Return(response, nil)

	dataFetcher.httpClientHandler = mock

	_, err := dataFetcher.FetchAllCharacters()

	assert.Error(t, err)
}

func TestFetchAllLocationsWhenOneElementIsReturned(t *testing.T) {
	response := model.HttpClientResponse{
		ResponseBody: []byte(oneLocationResponse),
		StatusCode:   http.StatusOK,
	}
	mock := new(httpClientMock)

	mock.On("Get", locationsPageOneUrl).Return(response, nil)

	dataFetcher.httpClientHandler = mock

	actual, err := dataFetcher.FetchAllLocations()

	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
	assert.Equal(t, oneElement, len(actual))
}

func TestFetchAllLocationssPagination(t *testing.T) {
	responsePageOne := model.HttpClientResponse{
		ResponseBody: []byte(paginationLocationResponse),
		StatusCode:   http.StatusOK,
	}
	responsePageTwo := model.HttpClientResponse{
		ResponseBody: []byte(oneLocationResponse),
		StatusCode:   http.StatusOK,
	}
	mock := new(httpClientMock)

	mock.On("Get", locationsPageOneUrl).Once().Return(responsePageOne, nil)
	mock.On("Get", locationsPageTwoUrl).Once().Return(responsePageTwo, nil)

	dataFetcher.httpClientHandler = mock

	actual, err := dataFetcher.FetchAllLocations()

	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
	assert.Equal(t, twoElements, len(actual))
}

func TestFetchAllLocationsWhenUnmarshalError(t *testing.T) {
	response := model.HttpClientResponse{
		ResponseBody: []byte(notAJSONResponseBody),
		StatusCode:   http.StatusOK,
	}
	mock := new(httpClientMock)

	mock.On("Get", charactersPageOneUrl).Return(response, nil)

	dataFetcher.httpClientHandler = mock

	_, err := dataFetcher.FetchAllCharacters()

	assert.Error(t, err)
}

func TestFetchAllCharactersWhenHttpClientError(t *testing.T) {
	httpClienErr := fmt.Errorf("httpClient error")
	mock := new(httpClientMock)

	mock.On("Get", charactersPageOneUrl).Return(model.HttpClientResponse{}, httpClienErr)

	dataFetcher.httpClientHandler = mock

	_, err := dataFetcher.FetchAllCharacters()

	assert.Error(t, err)
}
