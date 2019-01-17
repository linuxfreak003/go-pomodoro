package cmd

import (
	"github.com/linuxfreak003/go-pomodoro/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		profile := viper.GetString("profile")
		app := viper.GetString("app")
		host := viper.GetString("host")

		client.StartClient(profile, app, host, port)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().Uint16VarP(&port, "port", "p", 50051, "port server should bind to")
	viper.BindPFlag("port", clientCmd.Flags().Lookup("port"))

	clientCmd.Flags().String("app", "spotify", "music app to use")
	viper.BindPFlag("app", clientCmd.Flags().Lookup("app"))

	clientCmd.Flags().String("host", "127.0.0.1", "hostname")
	viper.BindPFlag("host", clientCmd.Flags().Lookup("host"))

}
