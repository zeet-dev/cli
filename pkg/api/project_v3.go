package api

import (
	"context"
)

// Shared fragment
var _ = `# @genqlient
		fragment ProjectV3AdapterSummary on ProjectV3Adapter {
			id
			name
		}`

func (c *Client) ListProjectV3s(ctx context.Context, filterInput FilterInput) (*ProjectV3AdapterConnection, error) {
	_ = `
		# @genqlient
		# @genqlient(for: "FilterInput.sort", pointer: true)
		# @genqlient(for: "FilterNode.criterion", pointer: true)
		# @genqlient(for: "FilterNode.expression", pointer: true)
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

		res, err := getProjectV3s(ctx, c.gql, user.Id, filterInput)
		adapters := &res.User.ProjectV3Adapters

		return adapters, err
}
