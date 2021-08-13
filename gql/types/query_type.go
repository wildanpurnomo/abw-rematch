package gqltypes

import "github.com/graphql-go/graphql"

var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user_contents": &graphql.Field{
				Type: &graphql.List{
					OfType: ContentType,
				},
				Description: "Get requesting user's content list",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	},
)
