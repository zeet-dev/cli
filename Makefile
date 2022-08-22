gen-docs:
	DOCS_OUT=./docs go run ./cmd/gen-docs
copy-docs:
	cp ./docs/* ../docs/docs/cli/commands
gen-go:
	get-graphql-schema http://localhost:7001/graphql > schema.graphql
	go generate ./...