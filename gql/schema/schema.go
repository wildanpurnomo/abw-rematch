package gqlschema

import (
	"github.com/graphql-go/graphql"
	gqlresolvers "github.com/wildanpurnomo/abw-rematch/gql/resolvers"
	gqltypes "github.com/wildanpurnomo/abw-rematch/gql/types"
)

func InitSchema() (graphql.Schema, error) {
	gqltypes.AddFieldToUserType("contents", &graphql.Field{
		Type: &graphql.List{
			OfType: gqltypes.ContentType,
		},
		Resolve: gqlresolvers.GetContentsByUserId,
	})

	gqltypes.AddFieldToContentType("author", &graphql.Field{
		Type:    gqltypes.UserType,
		Resolve: gqlresolvers.GetUserByIdResolver,
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    gqltypes.QueryType,
		Mutation: gqltypes.MutationType,
	}
	return graphql.NewSchema(schemaConfig)
}
