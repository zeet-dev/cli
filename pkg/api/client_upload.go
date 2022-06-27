package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type payload struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables,omitempty"`
}

type response struct {
	Data   interface{}   `json:"data"`
	Errors gqlerror.List `json:"errors"`
}

// Our graphql client doesn't support file uploads, so this function is a manual workaround which uses http
func uploadFile(client *http.Client, path string, query string, variables interface{}, fileVariable string, file []byte, retval interface{}) error {
	contentType, buf, err := makeForm(query, variables, fileVariable, file)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", path, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var dataAndErrors response
	dataAndErrors.Data = retval
	err = json.NewDecoder(resp.Body).Decode(&dataAndErrors)
	if err != nil {
		return err
	}

	if len(dataAndErrors.Errors) > 0 {
		return dataAndErrors.Errors
	}
	return nil
}

func makeForm(query string, variables interface{}, fileVariable string, file []byte) (string, *bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	defer w.Close()

	opPayload, err := json.Marshal(payload{Query: query, Variables: variables})
	if err != nil {
		return "", nil, err
	}
	_ = w.WriteField("operations", string(opPayload))
	_ = w.WriteField("map", fmt.Sprintf(`{"0": ["variables.%s"] }`, fileVariable))

	part, err := w.CreateFormFile("0", "file")
	if err != nil {
		return "", nil, err
	}
	_, err = part.Write(file)
	if err != nil {
		return "", nil, err
	}

	return w.FormDataContentType(), buf, nil
}
