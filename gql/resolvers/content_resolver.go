package gqlresolvers

import (
	"errors"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	gqldataloaders "github.com/wildanpurnomo/abw-rematch/gql/dataloaders"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

var (
	contentLoader     = dataloader.NewBatchedLoader(gqldataloaders.ContentBatchFn)
	DeleteContentById = func(params graphql.ResolveParams) (interface{}, error) {
		contextValue := libs.GetContextValues(params.Context)
		userId := contextValue.UserID
		if userId == "0" {
			return nil, errors.New("Invalid token or user not found")
		}

		contentId := params.Args["content_id"].(string)

		var content models.Content
		if err := repositories.Repo.GetContentByUserIdAndContentId(&content, userId, contentId); err != nil {
			return false, errors.New("Content not found or unauthorized user")
		}

		if err := repositories.Repo.DeleteContent(&content); err != nil {
			return false, errors.New("Whoops!")
		}

		return true, nil
	}
	GetContentsByUserId = func(params graphql.ResolveParams) (interface{}, error) {
		source := params.Source.(models.User)
		userId := fmt.Sprint(source.ID)

		thunk := contentLoader.Load(params.Context, dataloader.StringKey(userId))
		result, err := thunk()
		if err != nil {
			return nil, err
		}

		contentsMap := make(map[string][]models.Content)
		for _, item := range result.([]models.Content) {
			currentAuthorId := fmt.Sprint(item.UserID)
			if contentsMap[currentAuthorId] == nil {
				contentsMap[currentAuthorId] = make([]models.Content, 0)
			}
			contentsMap[currentAuthorId] = append(contentsMap[currentAuthorId], item)
		}

		return contentsMap[userId], nil
	}
	GetMyContentsResolver = func(params graphql.ResolveParams) (interface{}, error) {
		contextValue := libs.GetContextValues(params.Context)
		userId := contextValue.UserID
		if userId == "0" {
			return nil, errors.New("Invalid token or user not found")
		}

		var contents []models.Content
		if err := repositories.Repo.GetContentByUserId(&contents, userId); err != nil {
			return nil, errors.New("Whoops!")
		}

		return contents, nil
	}
)
