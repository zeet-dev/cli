gen-docs:
	DOCS_OUT=./docs go run ./cmd/gen-docs
copy-docs:
	cp ./docs/* ../docs/docs/cli/commands

build:
	go build -o ./bin/ ./cmd/zeet

link-dev:
	ln -s $(shell pwd)/bin/zeet ~/bin/zeet

gen-go:
	go generate ./...

# GraphQL

gen-gql: get-schema gen-gql-go gen-gql-sdk

get-schema-local:
	get-graphql-schema http://localhost:7001/graphql > schema_0.graphql
	get-graphql-schema http://localhost:7001/v1/graphql > schema_1.graphql

get-schema:
	get-graphql-schema https://anchor.zeet.co/graphql > schema_0.graphql
	get-graphql-schema https://anchor.zeet.co/v1/graphql > schema_1.graphql

gen-gql-go:
	go run github.com/Khan/genqlient genqlient.yaml

gen-gql-sdk:
	go run github.com/Khan/genqlient genqlient_0.yaml
	go run github.com/Khan/genqlient genqlient_1.yaml