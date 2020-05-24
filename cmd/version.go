package cmd

import (
	"github.com/linuxsuren/net-proxy/cmd/common"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates a version command
func NewVersionCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			version := common.GetVersion()
			cmd.Printf("Version: %s\n", version)
			cmd.Printf("Last Commit: %s\n", common.GetCommit())
			cmd.Printf("Build Date: %s\n", common.GetDate())
			return
		},
	}
	return
}
