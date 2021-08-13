package gqltypes

import "github.com/graphql-go/graphql"

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"profile_picture": &graphql.Field{
				Type: graphql.String,
			},
			"points": &graphql.Field{
				Type: graphql.Int,
			},
			"contents": &graphql.Field{
				Type: &graphql.List{
					OfType: ContentType,
				},
			},
		},
	},
)
