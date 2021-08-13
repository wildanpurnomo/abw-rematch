package gqlschema

import (
	"github.com/graphql-go/graphql"
	gqltypes "github.com/wildanpurnomo/abw-rematch/gql/types"
)

func InitSchema() (graphql.Schema, error) {
	schemaConfig := graphql.SchemaConfig{
		Query:    gqltypes.QueryType,
		Mutation: gqltypes.MutationType,
	}
	return graphql.NewSchema(schemaConfig)
}
