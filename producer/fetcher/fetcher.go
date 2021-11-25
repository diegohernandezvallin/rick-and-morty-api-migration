package fetcher

import (
	"encoding/json"
	"fmt"

	"github.com/rick-and-morty-character-migration/producer/httpclient"
	"github.com/rick-and-morty-character-migration/producer/model"
)

const (
	characterResource = "/character"
	locationResource  = "/location"

	pageParam = "page"
)

var dataFetcher DataFetcher

type Fetcher interface {
	FetchAllCharacters() ([]model.Character, error)
	FetchAllLocations() ([]model.Location, error)
}

type DataFetcher struct {
	httpClientHandler httpclient.HttpClient
	url               string
}

func NewDataFetcher(httpClientHandler httpclient.HttpClient, url string) DataFetcher {
	if dataFetcher.httpClientHandler == nil {
		dataFetcher = DataFetcher{
			httpClientHandler: httpClientHandler,
			url:               url,
		}
	}

	return dataFetcher
}

func (dataFetcher DataFetcher) FetchAllCharacters() ([]model.Character, error) {
	var results []model.Character

	url := fmt.Sprintf("%s%s?%s=1", dataFetcher.url, characterResource, pageParam)
	for url != "" {
		response, err := dataFetcher.fetch(url)
		if err != nil {
			return nil, fmt.Errorf("unable to get all characters. %v", err)
		}

		var characters model.CharacterResponse
		err = json.Unmarshal(response, &characters)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal response. %v", err)
		}

		results = append(results, characters.Results...)
		url = characters.Info.Next
	}

	return results, nil
}

func (dataFetcher DataFetcher) FetchAllLocations() ([]model.Location, error) {
	var results []model.Location

	url := fmt.Sprintf("%s%s?%s=1", dataFetcher.url, locationResource, pageParam)
	for url != "" {
		response, err := dataFetcher.fetch(url)
		if err != nil {
			return nil, fmt.Errorf("unable to get all locations. %v", err)
		}

		var locations model.LocationResponse
		err = json.Unmarshal(response, &locations)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal response. %v", err)
		}

		results = append(results, locations.Results...)
		url = locations.Info.Next
	}

	return results, nil
}

func (dataFetcher DataFetcher) fetch(url string) ([]byte, error) {
	httpClientResponse, err := dataFetcher.httpClientHandler.Get(url)
	if err != nil {
		return nil, err
	}

	if httpClientResponse.StatusCode < 200 || httpClientResponse.StatusCode > 299 {
		return nil, fmt.Errorf("failed request on url: %s. %v", url, err)
	}

	return httpClientResponse.ResponseBody, nil
}
