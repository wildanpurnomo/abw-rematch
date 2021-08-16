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
				Description: "Log out to system",
				Resolve:     gqlresolvers.LogoutResolver,
			},
			"update_username": &graphql.Field{
				Type:        UserType,
				Description: "Update requesting user's username",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: gqlresolvers.UpdateUsernameResolver,
			},
			"update_password": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Update requesting user's password",
				Args: graphql.FieldConfigArgument{
					"old_password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"new_password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: gqlresolvers.UpdatePasswordResolver,
			},
			"delete_content": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete one of requesting user's content",
				Args: graphql.FieldConfigArgument{
					"content_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: gqlresolvers.DeleteContentById,
			},
		},
	},
)
