package server

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/stuckinforloop/fabrik/server"
)

var port string

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, _ []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}

		srv := server.New(port)
		srv.Start()
	},
}

func init() {
	ServerCmd.Flags().StringVarP(&port, "port", "p", "", "Port to run the server on (default: 8080)")
}
