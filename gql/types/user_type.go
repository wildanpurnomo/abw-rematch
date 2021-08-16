package gqltypes

import (
	"github.com/graphql-go/graphql"
	gqlresolvers "github.com/wildanpurnomo/abw-rematch/gql/resolvers"
)

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "UserType",
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
				Resolve: gqlresolvers.GetContentsByUserId,
			},
		},
	},
)

func AddFieldToUserType(fieldName string, fieldConfig *graphql.Field) {
	UserType.AddFieldConfig(fieldName, fieldConfig)
}
