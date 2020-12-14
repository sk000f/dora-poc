package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/sk000f/metrix/pkg/http/graphql/graph/generated"
	"github.com/sk000f/metrix/pkg/http/graphql/graph/model"
)

func (r *queryResolver) AllProjectNames(ctx context.Context) ([]*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AllProjectGroupNames(ctx context.Context) ([]*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DeploymentFrequency(ctx context.Context, dateRange *model.DateRange, projectName *string, groupName *string) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ChangeFailRate(ctx context.Context, dateRange *model.DateRange, projectName *string, groupName *string) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MeanTimeToRecover(ctx context.Context, dateRange *model.DateRange, projectName *string, groupName *string) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ChangeLeadTime(ctx context.Context, dateRange *model.DateRange, projectName *string, groupName *string) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
