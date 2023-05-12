package reconciliation

import (
	gql "github.com/graphql-go/graphql"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-sdk/pkg/apis/graphql"
)

type reconcilerInformation struct {
	ID string `json:"id"`
}

type reconcilersApi struct {
}

func newReconcilersApi() *reconcilersApi {
	return &reconcilersApi{}
}

func (j *reconcilersApi) BindResource(ctx graphql.BindingContext) {
	reconcilerType := ctx.LookupOutputType(forddb.TypeOf(&reconcilerInformation{}))

	getJob := &gql.Field{
		Name: "Reconciler",
		Type: reconcilerType,

		Args: map[string]*gql.ArgumentConfig{
			"id": {
				Type: gql.NewNonNull(gql.String),
			},
		},

		Resolve: func(p gql.ResolveParams) (interface{}, error) {
			return &reconcilerInformation{}, nil
		},
	}

	getJobs := &gql.Field{
		Name: "allReconcilers",
		Type: gql.NewList(reconcilerType),

		Args: map[string]*gql.ArgumentConfig{
			"id":  {Type: gql.String},
			"ids": {Type: gql.NewList(gql.String)},
		},

		Resolve: func(p gql.ResolveParams) (interface{}, error) {
			return []*reconcilerInformation{}, nil
		},
	}

	ctx.RegisterQuery(getJob, getJobs)
}
