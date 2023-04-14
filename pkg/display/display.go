package display

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/zeet-dev/cli/pkg/api"
)

// a FetchFunction retrieves a page of data and returns it formatted as maps
type FetchFunction func(*api.PageInput) (*string, []table.Row, *api.PageInfo, int, error)

func DisplayPaginatedTable(fetch FetchFunction) (error) {
	reader := bufio.NewReader(os.Stdin)

	pageInput := &api.PageInput{
		First: 10,
		After: "0",
	}

	for ; ; {
		title, rows, pageInfo, totalCount, err := fetch(pageInput)
		if err != nil {
			return err
		}

		PrintTabularData(title, rows)

		start, err := strconv.Atoi(pageInfo.StartCursor)
		if err != nil {
			return err
		}
		fmt.Printf("Showing %d-%s of %d.\n", start + 1, pageInfo.EndCursor, totalCount)

		prompt := ""
		if pageInfo.HasNextPage {
			prompt = prompt + "[n]ext page | "
		}
		if pageInfo.HasPreviousPage {
			prompt = prompt + "[p]revious page | "
		}
		if prompt != "" {
			prompt = prompt + "[q]uit"
		} else {
			return nil
		}

		fmt.Println(prompt)

		data, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		action := strings.ToLower(strings.TrimSpace(data))
		switch action {
		case "q":
			return nil
		case "n":
			if pageInfo.HasNextPage {
				pageInput.After = pageInfo.EndCursor
			}
		case "p":
			if pageInfo.HasPreviousPage {
				pageInput.After = strconv.Itoa(start - pageInput.First)
			}
		}
	}
}

func PrintTabularData(title *string, rows []table.Row) {
	tw := table.NewWriter()
	tw.AppendHeader(rows[0])

	for i := 1; i < len(rows); i++ {
		tw.AppendRow(rows[i])
	}

	tw.SetTitle(*title)
	tw.SetStyle(table.StyleRounded)
	tw.Style().Options.SeparateRows = true
	fmt.Println(tw.Render())
}
