package api

import (
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/spf13/viper"
	"github.com/zeet-dev/pkg/utils"
	"k8s.io/client-go/transport"
)

func NewGraphQLClient() *graphql.Client {
	token := viper.GetString("auth.access_token")

	tp := http.DefaultTransport
	if viper.GetBool("debug") {
		tp = utils.LoggingHttpTransport
	}

	return graphql.NewClient(utils.URLJoin(viper.GetString("server"), "graphql"), &http.Client{
		Transport: transport.NewBearerAuthRoundTripper(token, tp),
	})
}
