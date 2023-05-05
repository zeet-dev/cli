gen-docs:
	DOCS_OUT=./docs go run ./cmd/gen-docs
copy-docs:
	cp ./docs/* ../docs/docs/cli/commands
gen-go:
	get-graphql-schema http://localhost:7001/graphql > schema.graphql
	go generate ./...
gen:
	go run github.com/Khan/genqlient

build:
	go build -o ./bin/ ./cmd/zeet

link-dev:
	ln -s $(shell pwd)/bin/zeet ~/bin/zeet
