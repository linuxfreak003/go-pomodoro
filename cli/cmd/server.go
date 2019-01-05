package cmd

import (
	pomo "github.com/linuxfreak003/go-pomodoro"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pomo.StartServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Uint16VarP(&port, "port", "p", 50051, "port server should bind to")
}
