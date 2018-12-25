package cmd

import (
	pomo "github.com/linuxfreak003/go-pomodoro"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pomo.StartClient()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
