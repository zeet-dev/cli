package build

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

func LintFiles(workDir, target string) error {
	if target == "lambda" {
		return LintFilesLambda(workDir)
	}
	return nil
}

func LintFilesLambda(workDir string) error {
	reqFile := filepath.Join(workDir, "requirements.txt")

	file, err := os.OpenFile(reqFile, os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	fmt.Println("Linting requirements.txt")

	bfile, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	file.Close()

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(bfile)
	if err != nil {
		return err
	}

	e, err := ianaindex.MIME.Encoding(result.Charset)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(transform.NewReader(bytes.NewBuffer(bfile), e.NewDecoder()))
	if err != nil {
		return err
	}

	out := ReplaceLambdaRequirements(string(content))
	if err := ioutil.WriteFile(reqFile, []byte(out), 0644); err != nil {
		return err
	}

	return nil
}
