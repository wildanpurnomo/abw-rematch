package gqlresolvers

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

var GetUserContentsResolver = func(params graphql.ResolveParams) (interface{}, error) {
	cookieAccess := libs.GetContextValues(params.Context)
	userId := cookieAccess.UserID
	if userId == "0" {
		return nil, errors.New("Invalid token or user not found")
	}

	var contents []models.Content
	if err := repositories.Repo.GetContentByUserId(&contents, userId); err != nil {
		return nil, errors.New("Whoops!")
	}

	return contents, nil
}
