package cmd

import (
	"github.com/spf13/cobra"
	"learn/rpc/grpc/gateway/server"
	"log"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC gateway server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recover error: %v\n", err)
			}
		}()

		server.Run()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&server.Port, "port", "p", "50051", "server port")
	serverCmd.Flags().StringVarP(&server.CertPemPath, "cert-pem", "", "./conf/certs/server.pem", "cert-pem path")
	serverCmd.Flags().StringVarP(&server.CertKeyPath, "cert-key", "", "./conf/certs/server.key", "cert-key path")
	serverCmd.Flags().StringVarP(&server.SererName, "cert-name", "", "grpc.abc", "server's hostname")
	rootCmd.AddCommand(serverCmd)
}
