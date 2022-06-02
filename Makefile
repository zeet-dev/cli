gen-docs:
	go run ./... gen-docs -d ./docs
copy-docs:
	cp ./docs/* ../docs/docs/cli/commands