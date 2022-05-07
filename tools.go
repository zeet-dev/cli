//go:build tools
// +build tools

// This package exists to cause `go mod` and `go get` to believe these tools
// are dependencies, even though they are not runtime dependencies of any
// package (these are tools used by our CI builds). This means they will appear
// in our `go.mod` file, but will not be a part of the build. Also, since the
// build target is something non-existent, these should not be included in any
// binaries.

package cli

import (
	_ "github.com/Khan/genqlient"
)
