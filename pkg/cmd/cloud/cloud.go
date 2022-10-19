package cloud

import (
	"github.com/spf13/cobra"
	"github.com/zeet-dev/cli/pkg/cmd/cloud/aws"
	"github.com/zeet-dev/cli/pkg/cmd/cloud/digitalocean"
	"github.com/zeet-dev/cli/pkg/cmd/cloud/gcp"
	"github.com/zeet-dev/cli/pkg/cmd/cloud/linode"
	"github.com/zeet-dev/cli/pkg/cmdutil"
)

func InitCloudCmds(f *cmdutil.Factory, root *cobra.Command) {
	awsCmd := &cobra.Command{
		Use:   "aws [command]",
		Short: "Manage AWS",
		Args:  cobra.ExactArgs(1),
	}

	awsCmd.AddCommand(aws.NewAWSLoginCmd(f))
	awsCmd.AddCommand(aws.NewAWSConsoleCmd(f))

	gcpCmd := &cobra.Command{
		Use:   "gcp [command]",
		Short: "Manage GCP",
		Args:  cobra.ExactArgs(1),
	}

	gcpCmd.AddCommand(gcp.NewGCPLoginCmd(f))

	linodeCmd := &cobra.Command{
		Use:   "linode [command]",
		Short: "Manage Linode",
		Args:  cobra.ExactArgs(1),
	}

	linodeCmd.AddCommand(linode.NewLinodeLoginCmd(f))

	doCmd := &cobra.Command{
		Use:   "do [command]",
		Short: "Manage DigitalOcean",
		Args:  cobra.ExactArgs(1),
	}

	doCmd.AddCommand(digitalocean.NewDOLoginCmd(f))

	root.AddCommand(awsCmd)
	root.AddCommand(gcpCmd)
	root.AddCommand(linodeCmd)
	root.AddCommand(doCmd)
}
