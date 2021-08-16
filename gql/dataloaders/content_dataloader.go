package gqldataloaders

import (
	"context"

	"github.com/graph-gophers/dataloader"
)

var ContentBatchFn = func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	return results
}
