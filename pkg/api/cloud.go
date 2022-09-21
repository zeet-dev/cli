package api

import (
	"context"

	"github.com/google/uuid"
)

func (c *Client) GetCloudAWS(ctx context.Context, cloudID uuid.UUID) (*GetCloudAWSResponse, error) {
	_ = `# @genqlient
		query GetCloudAWS($id: UUID!) {
		  currentUser {
			awsAccount(id: $id) {
			  id
			  roleARN
			  externalID
			}
		  }
		}
	`

	return GetCloudAWS(ctx, c.gql, cloudID)
}

func (c *Client) GetCloudGCP(ctx context.Context, cloudID uuid.UUID) (*GetCloudGCPResponse, error) {
	_ = `# @genqlient
		query GetCloudGCP($id: UUID!) {
		  currentUser {
			gcpAccount(id: $id) {
			  id
			  credentials
			}
		  }
		}
	`

	return GetCloudGCP(ctx, c.gql, cloudID)
}

func (c *Client) GetCloudLinode(ctx context.Context, cloudID uuid.UUID) (*GetCloudLinodeResponse, error) {
	_ = `# @genqlient
		query GetCloudLinode($id: UUID!) {
		  currentUser {
			linodeAccount(id: $id) {
			  id
			  accessToken
			}
		  }
		}
	`

	return GetCloudLinode(ctx, c.gql, cloudID)
}

func (c *Client) GetCloudDigitalOcean(ctx context.Context, cloudID uuid.UUID) (*GetCloudDigitalOceanResponse, error) {
	_ = `# @genqlient
		query GetCloudDigitalOcean($id: UUID!) {
		  currentUser {
			doAccount(id: $id) {
			  id
			  accessToken
			}
		  }
		}
	`

	return GetCloudDigitalOcean(ctx, c.gql, cloudID)
}
