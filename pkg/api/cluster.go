package api

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type Cluster struct {
	Id uuid.UUID `copier:"Id"`
}

func (c *Client) UpdateClusterKubeconfig(ctx context.Context, clusterId uuid.UUID, kubeconfig []byte) (*Cluster, error) {
	out := &Cluster{}

	query := `# @genqlient
		mutation updateCluster($id: UUID!, $file: Upload!) {
		  updateCluster(input: {id: $id, kubeconfig: $file}) {
			id
		  }
		}
	`

	type input struct {
		File []byte    `json:"file"`
		Id   uuid.UUID `json:"id"`
	}
	var res updateClusterResponse
	err := c.upload.UploadFile(query, input{File: kubeconfig, Id: clusterId}, "file", kubeconfig, res)
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(&out, res.UpdateCluster); err != nil {
		return out, err
	}

	return out, nil
}

type GetClusterKubeconfigResponse struct {
	ID         uuid.UUID
	Name       string
	Kubeconfig string
}

func (c *Client) GetClusterKubeconfig(ctx context.Context, clusterID uuid.UUID) (*GetClusterKubeconfigResponse, error) {
	_ = `# @genqlient
		query getClusterKubeconfig($id: UUID!) {
		  currentUser {
			cluster(id: $id) {
			  id
			  name
			  kubeconfig
			}
		  }
		}
	`

	res, err := GetClusterKubeconfigQuery(ctx, c.gql, clusterID)
	if err != nil {
		return nil, err
	}

	return &GetClusterKubeconfigResponse{
		ID:         res.CurrentUser.Cluster.Id,
		Name:       res.CurrentUser.Cluster.Name,
		Kubeconfig: res.CurrentUser.Cluster.Kubeconfig,
	}, nil
}

type ListClustersResponse struct {
	ID              uuid.UUID
	Name            string
	CloudProvider   string
	ClusterProvider string
	Region          string
	Connected       bool
}

func (c *Client) ListClusters(ctx context.Context, path string) ([]ListClustersResponse, error) {
	_ = `# @genqlient
		query listClusters {
		  currentUser {
			clusters {
			  id
			  name
			  cloudProvider
			  clusterProvider
			  region
			  connected
			}
		  }
		}
	`

	_ = `# @genqlient
	query listClustersForTeam($path: String) {
		team(path: $path) {
			user {
				clusters {
					id
					name
					cloudProvider
					clusterProvider
					region
					connected
				}
			}
	    }
	}
	`

	if path != "" {
		res, err := ListClustersForTeamQuery(ctx, c.gql, path)
		if err != nil {
			if len(res.Team.User.Clusters) > 0 {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
			} else {
				return nil, err
			}
		}

		out := make([]ListClustersResponse, len(res.Team.User.Clusters))
		for i, cluster := range res.Team.User.Clusters {
			out[i] = ListClustersResponse{
				ID:              cluster.Id,
				Name:            cluster.Name,
				CloudProvider:   string(cluster.CloudProvider),
				ClusterProvider: string(cluster.ClusterProvider),
				Region:          cluster.Region,
				Connected:       cluster.Connected,
			}
		}

		return out, nil
	}

	res, err := ListClustersQuery(ctx, c.gql)
	if err != nil {
		if len(res.CurrentUser.Clusters) > 0 {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
		} else {
			return nil, err
		}
	}

	out := make([]ListClustersResponse, len(res.CurrentUser.Clusters))
	for i, cluster := range res.CurrentUser.Clusters {
		out[i] = ListClustersResponse{
			ID:              cluster.Id,
			Name:            cluster.Name,
			CloudProvider:   string(cluster.CloudProvider),
			ClusterProvider: string(cluster.ClusterProvider),
			Region:          cluster.Region,
			Connected:       cluster.Connected,
		}
	}

	return out, nil
}
