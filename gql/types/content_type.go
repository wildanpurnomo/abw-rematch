package gqltypes

import "github.com/graphql-go/graphql"

var ContentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Content",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"body": &graphql.Field{
				Type: graphql.String,
			},
			"media_urls": &graphql.Field{
				Type: &graphql.List{
					OfType: graphql.String,
				},
			},
			"youtube_url": &graphql.Field{
				Type: graphql.String,
			},
			"slug": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
