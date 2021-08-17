package gqldataloaders

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

var ContentBatchFn = func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	var result dataloader.Result

	var contents []models.Content
	if err := repositories.Repo.GetContentInUserIds(&contents, keys.Keys()); err != nil {
		result.Error = err
	} else {
		result.Data = contents
	}

	results = append(results, &result)
	return results
}
