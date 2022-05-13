package cmd

import (
	"fmt"
	"time"

	"github.com/zeet-dev/cli/pkg/api"
	"github.com/zeet-dev/cli/pkg/utils"
)

func printOrPollLogs[S comparable](getLogs func() ([]api.LogEntry, error), getStatus func() (S, error), poll bool) error {
	if poll {
		if err := pollLogs(getLogs, getStatus); err != nil {
			return err
		}
	} else {
		logs, err := getLogs()
		if err != nil {
			return err
		}
		for _, log := range logs {
			fmt.Println(log.Text)
		}
	}

	return nil
}

// pollLogs continuously prints the output of getLogs, and stops when the output of getStatus changes
func pollLogs[S comparable](getLogs func() ([]api.LogEntry, error), getStatus func() (S, error)) (err error) {
	lastLog := 0

	initialStatus, err := getStatus()
	if err != nil {
		return err
	}

	shouldContinue := true

	for shouldContinue {
		// Stop looping if the status changes
		status, err := getStatus()
		if err != nil {
			return err
		}
		shouldContinue = status == initialStatus

		logs, err := getLogs()
		if err != nil {
			return err
		}
		// Catch the edge case where the deployment has been cancelled after getStatus was called
		// but before getLogs, making getLogs return [] and the range panicking
		if len(logs) == 0 && lastLog > 0 {
			return nil
		}

		// TODO
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
