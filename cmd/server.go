package cmd

import (
	"github.com/linuxfreak003/go-pomodoro/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer(port, token, channel)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Uint16VarP(&port, "port", "p", 50051, "port server should bind to")
	serverCmd.Flags().StringVarP(&token, "token", "t", "", "slack Legacy API token")
	serverCmd.Flags().StringVarP(&channel, "channel", "c", "pomodoro-spotify", "slack channel to send message to")
}
