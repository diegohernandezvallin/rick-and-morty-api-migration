package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rick-and-morty-character-migration/producer/fetcher"
	"github.com/rick-and-morty-character-migration/producer/httpclient"
	"github.com/rick-and-morty-character-migration/producer/model"
	"github.com/rick-and-morty-character-migration/producer/publishing"
	"github.com/rick-and-morty-character-migration/producer/publishing/kafka"
)

func main() {
	var (
		brokers = os.Getenv("KAFKA_BROKERS")
		topic   = os.Getenv("KAFKA_TOPIC")
		url     = os.Getenv("RICK_AND_MORTY_URL")

		dataFetcher fetcher.Fetcher
	)

	log.Printf("KAFKA_BROKERS: %s\nKAFKA_TOPIC: %s\nURL: %s\n", brokers, topic, url)
	client := http.Client{}
	httpClientHandler := httpclient.NewHttpHandler(&client)
	dataFetcher = fetcher.NewDataFetcher(httpClientHandler, url)

	publisher := kafka.NewPublisher(strings.Split(brokers, ","), topic)

	r := gin.Default()
	r.POST("/migrate/locations", migrateLocations(publisher, dataFetcher))
	r.POST("/migrate/characters", migrateCharacters(publisher, dataFetcher))

	_ = r.Run()
}

func migrateLocations(publisher publishing.Publisher, dataFetcher fetcher.Fetcher) func(*gin.Context) {
	return func(c *gin.Context) {
		locations, err := dataFetcher.FetchAllLocations()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		var notPublishedLocations []string
		for _, location := range locations {
			id := strconv.Itoa(location.ID)
			headers := map[string]string{"Type": "location", "SourceId": id}
			message := model.Message{
				Payload: location,
				Headers: headers,
				Key:     id,
			}
			err := publisher.Publish(context.Background(), message)
			if err != nil {
				errMessage := fmt.Sprintf("location with id: %s not published to kafka. error: %v", id, err)
				notPublishedLocations = append(notPublishedLocations, errMessage)
			}
		}

		if len(notPublishedLocations) != 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": notPublishedLocations})
		} else {
			c.JSON(http.StatusAccepted, gin.H{"message": "Locations published"})
		}
	}
}

func migrateCharacters(publisher publishing.Publisher, dataFetcher fetcher.Fetcher) func(*gin.Context) {
	return func(c *gin.Context) {
		characters, err := dataFetcher.FetchAllCharacters()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		var notPublishedCharacters []string
		for _, location := range characters {
			id := strconv.Itoa(location.ID)
			headers := map[string]string{"Type": "location", "SourceId": id}
			message := model.Message{
				Payload: location,
				Headers: headers,
				Key:     id,
			}
			err := publisher.Publish(context.Background(), message)
			if err != nil {
				errMessage := fmt.Sprintf("character with id: %s not published to kafka", id)
				notPublishedCharacters = append(notPublishedCharacters, errMessage)
			}
		}

		if len(notPublishedCharacters) != 0 {
			for _, notPublishedCharacter := range notPublishedCharacters {
				log.Println(notPublishedCharacter)
			}
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "Locations published"})
	}
}
