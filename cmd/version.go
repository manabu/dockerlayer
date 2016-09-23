package cmd

import (
	"fmt"

	"github.com/manabu/dockerlayer/config"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dockerlayer",
	Long:  `Print the version number of dockerlayer`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dockerlayer version %s (%s)\n", config.VersionString, config.CommitID)
	},
}
