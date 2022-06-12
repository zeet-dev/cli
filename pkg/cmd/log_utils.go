package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/utils"
)

func pollLogs(getLogs func() ([]api.LogEntry, error), shouldContinue func() (bool, error), out io.Writer) error {
	lastLog := 0

	cont, err := shouldContinue()
	if err != nil {
		return err
	}
	for cont {
		cont, err = shouldContinue()
		if err != nil {
			return err
		}

		logs, err := getLogs()
		if err != nil {
			return err
		}
		// Catch the edge case where the deployment has been cancelled after shouldContinue was called
		// but before getLogs, making getLogs return [] and the range panicking
		if len(logs) == 0 && lastLog > 0 {
			return nil
		}

		if len(logs) < lastLog {
			return nil
		}

		// Sometimes the backend returns an empty log which will then be replaced (same index) the next request...
		logs = utils.SliceFilter(logs, func(l api.LogEntry) bool {
			return l.Text != ""
		})

		for _, log := range logs[lastLog:] {
			fmt.Println(log.Text)
		}
		lastLog = len(logs)

		time.Sleep(time.Second)
	}

	return nil
}
