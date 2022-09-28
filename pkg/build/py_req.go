package build

import (
	"fmt"
	"regexp"
	"strings"
)

type PythonRequirement struct {
	Name    string
	Version string
}

var (
	lambdaRequirementReplacement = map[string]PythonRequirement{
		"psycopg2": {
			Name:    "psycopg2-binary",
			Version: "2.9.3",
		},
		"pysqlite": {
			Name:    "pysqlite3-binary",
			Version: "0.4.6",
		},
	}
	lambdaReplacedMap = map[string]bool{
		"psycopg2-binary":  true,
		"pysqlite3-binary": true,
	}
)

func ReplaceLambdaRequirements(reqFile string) string {
	replacedReqs := map[string]bool{}

	out := []string{}
	for _, line := range regexp.MustCompile("\r?\n").Split(reqFile, -1) {
		re := regexp.MustCompile("[ =><]+")
		splits := re.Split(line, 2)
		pylib := strings.ToLower(splits[0])

		if _, ok := replacedReqs[pylib]; ok {
			continue
		}

		if repl, ok := lambdaRequirementReplacement[pylib]; ok {
			replacedReqs[repl.Name] = true
			out = append(out, fmt.Sprintf("%s==%s", repl.Name, repl.Version))
		} else if _, ok := lambdaReplacedMap[pylib]; ok {
			replacedReqs[pylib] = true
			out = append(out, line)
		} else {
			out = append(out, line)
		}
	}

	return strings.Join(out, "\n")
}
