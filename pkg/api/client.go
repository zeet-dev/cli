package api

import (
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	zutils "github.com/zeet-dev/pkg/utils"
	"k8s.io/client-go/transport"
)

type Client struct {
	GQL graphql.Client
}

func New(server, token, version string, debug bool) *Client {
	client := newGraphQLClient(server, token, version, debug)

	return &Client{GQL: client}
}

func newGraphQLClient(server, token, version string, debug bool) graphql.Client {
	tp := http.DefaultTransport
	if debug {
		tp = zutils.LoggingHttpTransport
	}
	tp = NewUserAgentTransport(fmt.Sprintf("zeet-cli/%s", version), tp)

	return graphql.NewClient(zutils.URLJoin(server, "graphql"), &http.Client{
		Transport: transport.NewBearerAuthRoundTripper(token, tp),
	})
}
