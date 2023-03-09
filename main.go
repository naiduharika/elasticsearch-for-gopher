package main

import (
	"context"
	"es4gophers/logic"
)

func main() {
	// context to share data between functions
	ctx := context.Background()

	ctx = logic.LoadMoviesFromFile(ctx)
	ctx = logic.ConnectWithElasticsearch(ctx)
	// logic.IndexMoviesAsDocuments(ctx)
	logic.QueryMovieByDocumentID(ctx)
	logic.BestKeanuActionMovies(ctx)
	logic.MovieCountPerGenreAgg(ctx)
	logic.BestKeanuActionMoviesAsync(ctx)
}