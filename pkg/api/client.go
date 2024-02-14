package api

import (
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	zutils "github.com/zeet-dev/pkg/utils"
	"k8s.io/client-go/transport"
)

type Client struct {
	gql    graphql.Client
	upload *UploadClient

	gqlV1    graphql.Client
	uploadV1 *UploadClient
}

func New(server, token, version string, debug bool) *Client {
	httpClient := newHTTPClient(debug, version, token)

	path := zutils.URLJoin(server, "graphql")
	pathV1 := zutils.URLJoin(server, "v1/graphql")

	return &Client{
		gql:    newGraphQLClient(httpClient, path),
		upload: NewUploadClient(httpClient, path),

		gqlV1:    newGraphQLClient(httpClient, pathV1),
		uploadV1: NewUploadClient(httpClient, pathV1),
	}
}

func newHTTPClient(debug bool, version, token string) *http.Client {
	tp := http.DefaultTransport
	if debug {
		tp = zutils.LoggingHttpTransport
	}
	tp = transport.NewUserAgentRoundTripper(fmt.Sprintf("zeet-cli/%s", version), tp)

	return &http.Client{
		Transport: transport.NewBearerAuthRoundTripper(token, tp),
	}
}

func newGraphQLClient(httpClient *http.Client, path string) graphql.Client {
	return graphql.NewClient(path, httpClient)
}

func (c *Client) Client() graphql.Client {
	return c.gql
}

func (c *Client) UploadClient() *UploadClient {
	return c.upload
}

func (c *Client) ClientV1() graphql.Client {
	return c.gqlV1
}

func (c *Client) UploadClientV1() *UploadClient {
	return c.uploadV1
}
