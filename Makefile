gen-docs:
	DOCS_OUT=./docs go run ./cmd/gen-docs
copy-docs:
	cp ./docs/* ../docs/docs/cli/commands