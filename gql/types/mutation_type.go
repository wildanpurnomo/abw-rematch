package gqltypes

import (
	"github.com/graphql-go/graphql"
	gqlresolvers "github.com/wildanpurnomo/abw-rematch/gql/resolvers"
)

var AuthFieldConfigArguments = graphql.FieldConfigArgument{
	"username": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"password": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
}

var MutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Type:        UserType,
				Description: "Register new user",
				Args:        AuthFieldConfigArguments,
				Resolve:     gqlresolvers.RegisterResolver,
			},
			"login": &graphql.Field{
				Type:        UserType,
				Description: "Log in to system",
				Args:        AuthFieldConfigArguments,
				Resolve:     gqlresolvers.LoginResolver,
			},
			"logout": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Log out to system and clear the jwt cookie",
				Resolve:     gqlresolvers.LogoutResolver,
			},
		},
	},
)
