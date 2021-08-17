package gqltypes

import (
	"github.com/graphql-go/graphql"
	gqlresolvers "github.com/wildanpurnomo/abw-rematch/gql/resolvers"
)

var (
	queryArgs = graphql.FieldConfigArgument{
		"limit": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"offset": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	}
	QueryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"authenticate": &graphql.Field{
					Type:        UserType,
					Description: "Get requesting user's data",
					Resolve:     gqlresolvers.AuthenticateResolver,
				},
				"my_contents": &graphql.Field{
					Type: &graphql.List{
						OfType: ContentType,
					},
					Args:        queryArgs,
					Description: "Get requesting user's content list",
					Resolve:     gqlresolvers.GetMyContentsResolver,
				},
			},
		},
	)
)
