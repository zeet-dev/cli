package cmdutil

import "github.com/zeet-dev/cli/internal/config"

func CheckAuth(c config.Config) bool {
	return c.GetString("auth.access_token") != ""
}
