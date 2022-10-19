package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/cmdutil"
	"k8s.io/utils/pointer"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func NewAWSConsoleCmd(f *cmdutil.Factory) *cobra.Command {
	var opts = &AWSLoginOptions{}
	opts.IO = f.IOStreams
	opts.ApiClient = f.ApiClient

	cmd := &cobra.Command{
		Use:  "console <cloud id>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := uuid.Parse(args[0])
			if err != nil {
				return errors.New("invalid cloud ID format")
			}
			opts.CloudID = id

			if err := runAWSConsole(opts); err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "error: %s", err.Error())
				return err
			}

			return nil
		},
	}

	return cmd
}

func runAWSConsole(opts *AWSLoginOptions) error {
	client, err := opts.ApiClient()
	if err != nil {
		return err
	}

	cloud, err := client.GetCloudAWS(context.Background(), opts.CloudID)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	stsClient := sts.NewFromConfig(cfg)

	input := &sts.AssumeRoleInput{
		RoleArn:         &cloud.CurrentUser.AwsAccount.RoleARN,
		ExternalId:      &cloud.CurrentUser.AwsAccount.ExternalID,
		RoleSessionName: pointer.StringPtr("zeet"),
		DurationSeconds: pointer.Int32Ptr(3600),
	}

	role, err := stsClient.AssumeRole(ctx, input)
	if err != nil {
		return err
	}

	sessionBytes, err := json.Marshal(map[string]*string{
		"sessionId":    role.Credentials.AccessKeyId,
		"sessionKey":   role.Credentials.SecretAccessKey,
		"sessionToken": role.Credentials.SessionToken,
	})
	if err != nil {
		return err
	}

	tokenURL, err := url.Parse("https://signin.aws.amazon.com/federation")
	if err != nil {
		return err
	}
	tokenURL.RawQuery = url.Values{
		"Action":          []string{"getSigninToken"},
		"Session":         []string{string(sessionBytes)},
		"SessionDuration": []string{"3500"},
	}.Encode()

	tokenResp, err := http.Get(tokenURL.String())
	if err != nil {
		return err
	}

	tokenBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		fmt.Fprintf(opts.IO.ErrOut, "failed to get signin token: %s", string(tokenBody))
		return err
	}

	var tokenObj map[string]string
	if err := json.Unmarshal(tokenBody, &tokenObj); err != nil {
		return err
	}

	signinURL, err := url.Parse("https://signin.aws.amazon.com/federation")
	if err != nil {
		return err
	}
	signinURL.RawQuery = url.Values{
		"Action":      []string{"login"},
		"SigninToken": []string{tokenObj["SigninToken"]},
		"Destination": []string{"https://console.aws.amazon.com/"},
		"Issuer":      []string{"https://zeet.co"},
	}.Encode()

	if !openBrowser(signinURL.String()) {
		fmt.Fprintf(opts.IO.Out, "Open the following URL in your browser:\n%s", signinURL.String())
	}

	return nil
}

func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
