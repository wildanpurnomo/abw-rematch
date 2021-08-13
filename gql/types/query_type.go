package gqltypes

import (
	"github.com/graphql-go/graphql"
	gqlresolvers "github.com/wildanpurnomo/abw-rematch/gql/resolvers"
)

var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"authenticate": &graphql.Field{
				Type:        UserType,
				Description: "Get requesting user's data",
				Resolve:     gqlresolvers.AuthenticateResolver,
			},
			"user_contents": &graphql.Field{
				Type: &graphql.List{
					OfType: ContentType,
				},
				Description: "Get requesting user's content list",
				Resolve:     gqlresolvers.GetUserContentsResolver,
			},
		},
	},
)
