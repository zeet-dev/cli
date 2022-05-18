package update

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rogpeppe/go-internal/semver"
	"github.com/zeet-dev/cli/internal/config"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	ID      int
	URL     string `json:"html_url"`
}

func Check(client *http.Client, state config.Config, currentVersion string) (*ReleaseInfo, error) {
	defer func() {
		state.Set("last_update_check", time.Now().Unix())
		state.WriteConfig()
	}()

	rel := &ReleaseInfo{}

	res, err := client.Get("https://api.github.com/repos/zeet-dev/cli/releases/latest")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, rel); err != nil {
		return nil, err
	}

	if c := semver.Compare(currentVersion, rel.TagName); c == -1 {
		return rel, nil
	} else {
		return nil, nil
	}
}

func ShouldCheck(state config.Config) bool {
	t := time.Unix(state.GetInt64("last_update_check"), 0)
	return time.Since(t).Hours() > 24
}
