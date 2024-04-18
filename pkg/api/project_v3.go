package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// Shared fragment
var _ = `# @genqlient
		fragment ProjectV3AdapterSummary on ProjectV3Adapter {
			id
			name
			projectV3 {
				id
				name
			}
			repo {
				id
				name
			}
		}`

func (c *Client) ListProjectV3s(ctx context.Context, filterInput *FilterInput) (*ProjectV3AdapterConnection, error) {
	_ = `
		# @genqlient(pointer: true)
		# @genqlient(for: "PageInput.first", pointer: false)
		# @genqlient(for: "PageInput.after", pointer: false)
		# @genqlient(for: "PageInfo.startCursor", pointer: false)
		# @genqlient(for: "PageInfo.endCursor", pointer: false)
		# @genqlient(for: "PageInfo.hasNextPage", pointer: false)
		# @genqlient(for: "PageInfo.hasPreviousPage", pointer: false)
		query getProjectV3s(
			$userId: ID!, $filter: FilterInput!) {
			user(id: $userId) {
				# @genqlient(typename: "ProjectV3AdapterConnection")
				projectV3Adapters(
					filter: $filter
				) {
					totalCount
					# @genqlient(typename: "PageInfo")
					pageInfo {
						startCursor
						endCursor
						hasNextPage
						hasPreviousPage
					}
					nodes {
						...ProjectV3AdapterSummary
					}
				}
			}
		}`

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	userID := user.Id.String()
	res, err := GetProjectV3sQuery(ctx, c.gql, &userID, filterInput)
	adapters := res.User.ProjectV3Adapters

	return adapters, err
}

func (c *Client) GetProjectV3(ctx context.Context, projectId uuid.UUID) (*ProjectV3AdapterSummary, error) {
	filter := &FilterInput{
		Page: &PageInput{
			First: 1,
			After: "0",
		},
		Filter: &FilterNode{
			Criterion: &FilterCriterion{
				ResourceAdapterFilter: &ResourceAdapterFilter{
					Ids: &MultiEntityCriterion{
						Value: []*uuid.UUID{
							&projectId,
						},
					},
				},
			},
		},
	}

	conn, err := c.ListProjectV3s(ctx, filter)
	if err != nil {
		return nil, err
	}

	adapter := &ProjectV3AdapterSummary{}
	err = copier.Copy(adapter, conn.Nodes[0])
	return adapter, err
}
