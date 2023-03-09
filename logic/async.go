package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"es4gophers/domain"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func BestKeanuActionMoviesAsync(ctx context.Context) {

	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)

	var searchBuffer bytes.Buffer
	search := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"Year": map[string]int{
								"gte": 1900,
								"lte": 2021,
							},
						},
					},
				},
			},
		},
	}
	err := json.NewEncoder(&searchBuffer).Encode(search)
	if err != nil {
		panic(err)
	}

	var numberOfDocs int = 0
	asyncRequest := esapi.AsyncSearchSubmitRequest{
		Index:          []string{"movies"},
		Body:           &searchBuffer,
		TrackTotalHits: true,
		Pretty:         true,
		Size:           &numberOfDocs,
	}
	response, err := asyncRequest.Do(ctx, client.Transport)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	time.Sleep(8 * time.Second)

	var SearchResponse = domain.SearchResponse{}
	err = json.NewDecoder(response.Body).Decode(&SearchResponse)
	if err != nil {
		panic(err)
	}

	if SearchResponse.Hits.Total.Value > 0 {
		var movieTitles []string
		for _, movieTitle := range SearchResponse.Hits.Hits {
			movieTitles = append(movieTitles, movieTitle.Source.Title)
		}
		fmt.Printf("✅ Best action movies from Keanu async: [%s] \n", strings.Join(movieTitles, ", "))
	}
}
