package gqltypes

import (
	"github.com/graphql-go/graphql"
	gqlresolvers "github.com/wildanpurnomo/abw-rematch/gql/resolvers"
)

var MutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Type:        UserType,
				Description: "Register new user",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: gqlresolvers.RegisterResolver,
			},
		},
	},
)
