package gqldataloaders

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

var UserBatchFn = func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	var result dataloader.Result

	var users []models.User
	if err := repositories.Repo.FetchUserInUserIds(&users, keys.Keys()); err != nil {
		result.Error = err
	} else {
		result.Data = users
	}

	results = append(results, &result)
	return results
}
