package api

import (
	"context"

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
	err := uploadFile(c.http, c.path, query, input{File: kubeconfig, Id: clusterId}, "file", kubeconfig, res)
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(&out, res.UpdateCluster); err != nil {
		return out, err
	}

	return out, nil
}
