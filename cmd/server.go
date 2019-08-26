package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
	var slackToken string

	home, _ := os.UserHomeDir()
	tokenFile, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/go-pomodoro/slack-token", home))

	if err == nil {
		slackToken = string(tokenFile)
	}

	if slackToken == "" {
		slackToken = os.Getenv("SLACK_TOKEN")
	}

	slackToken = strings.TrimSpace(slackToken)

	serverCmd.Flags().Uint16VarP(&port, "port", "p", 1337, "port server should bind to")
	serverCmd.Flags().StringVarP(&token, "token", "t", slackToken, "slack Legacy API token")
	serverCmd.Flags().StringVarP(&channel, "channel", "c", "pomodoro-spotify", "slack channel to send message to")
}
