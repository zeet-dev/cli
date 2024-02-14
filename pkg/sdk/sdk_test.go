package sdk_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
	zeetv0 "github.com/zeet-dev/cli/pkg/sdk/v0"
	zeetv1 "github.com/zeet-dev/cli/pkg/sdk/v1"
)

func TestSDKV0(t *testing.T) {
	httpClient := http.Client{}
	graphqlClient := graphql.NewClient("https://anchor.zeet.co/graphql", &httpClient)

	user, err := zeetv0.CurrentUserQuery(context.Background(), graphqlClient)

	// TODO: test with mock server
	require.ErrorContains(t, err, "input: currentUser not authenticated")
	require.Empty(t, user)
}

func TestSDKV1(t *testing.T) {
	httpClient := http.Client{}
	graphqlClient := graphql.NewClient("https://anchor.zeet.co/v1/graphql", &httpClient)

	user, err := zeetv1.CurrentUserQuery(context.Background(), graphqlClient)

	// TODO: test with mock server
	require.ErrorContains(t, err, "input: currentUser not authenticated")
	require.Empty(t, user)
}
