package api

import (
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/spf13/viper"
	"github.com/zeet-co/anchor/pkg/utils"
	"k8s.io/client-go/transport"
)

func NewGraphQLClient() *graphql.Client {
	token := viper.GetString("auth.access_token")

	return graphql.NewClient(utils.URLJoin(viper.GetString("server"), "graphql"), &http.Client{
		Transport: transport.NewBearerAuthRoundTripper(token, http.DefaultTransport),
	})
}
