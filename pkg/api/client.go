package api

import (
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	zutils "github.com/zeet-dev/pkg/utils"
	"k8s.io/client-go/transport"
)

type Client struct {
	gql  graphql.Client
	http *http.Client
	path string
}

func New(server, token, version string, debug bool) *Client {
	path := zutils.URLJoin(server, "graphql")
	httpClient := newHTTPClient(debug, version, token)
	gqlClient := newGraphQLClient(path, httpClient)

	return &Client{path: path, gql: gqlClient, http: httpClient}
}

func newHTTPClient(debug bool, version, token string) *http.Client {
	tp := http.DefaultTransport
	if debug {
		tp = zutils.LoggingHttpTransport
	}
	tp = NewUserAgentTransport(fmt.Sprintf("zeet-cli/%s", version), tp)

	return &http.Client{
		Transport: transport.NewBearerAuthRoundTripper(token, tp),
	}
}

func newGraphQLClient(path string, httpClient *http.Client) graphql.Client {
	return graphql.NewClient(path, httpClient)
}
