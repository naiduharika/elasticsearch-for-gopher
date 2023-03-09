package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"es4gophers/domain"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func IndexMoviesAsDocuments(ctx context.Context) {

	movies := ctx.Value(domain.MoviesKey).([]domain.Movie)
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)

	// for documentID, document := range movies {
	// 	res, err := client.Index("movies", esutil.NewJSONReader(document),
	// 				client.Index.WithDocumentID(strconv.Itoa(documentID)))

	// 	if err == nil {
	// 		fmt.Println(res)
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	// }

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      "movies",
		Client:     client,
		NumWorkers: 5,
	})
	if err != nil {
		panic(err)
	}

	var countSuccessful uint64

	for documentID, document := range movies {
		data, err := json.Marshal(&document)
		if err != nil {
			panic(err)
		}
		err = bulkIndexer.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(documentID),
				Body:       bytes.NewReader(data),

				OnSuccess: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						fmt.Printf("ERROR: %d", err)
					} else {
						fmt.Printf("ERROR %s: %s", biri.Error.Type, biri.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			panic(err)
		}
	}

	bulkIndexer.Close(ctx)
	biStats := bulkIndexer.Stats()
	fmt.Printf("✅ Movies indexed on Elasticsearch: %d \n", biStats.NumIndexed)
}
