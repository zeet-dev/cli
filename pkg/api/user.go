package api

import (
	"context"
	"fmt"

	graphql "github.com/hasura/go-graphql-client"
)

func ShowCurrentUser(ctx context.Context) error {
	client := NewGraphQLClient()

	var query struct {
		CurrentUser struct {
			ID    graphql.String
			Login graphql.String
		} `graphql:"currentUser()"`
	}

	err := client.Query(ctx, &query, map[string]interface{}{})
	if err != nil {
		return err
	}

	fmt.Println("You are logged in as:", query.CurrentUser.Login)
	return nil
}
