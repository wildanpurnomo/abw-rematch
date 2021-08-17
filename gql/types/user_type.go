package gqltypes

import (
	"github.com/graphql-go/graphql"
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
		},
	},
)

func AddFieldToUserType(fieldName string, fieldConfig *graphql.Field) {
	UserType.AddFieldConfig(fieldName, fieldConfig)
}
