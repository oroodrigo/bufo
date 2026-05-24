package cli

import (
	"fmt"

	"github.com/oroodrigo/bufo/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects in the proxy",
	Run: func(cmd *cobra.Command, args []string) {
		socketFile := config.Load().SocketFile

		if err := ensureDaemon(socketFile); err != nil {
			fmt.Println("Erro ao garantir que o daemon está em execução:", err)
			return
		}

		client := NewApiClient(socketFile)
		routes, err := client.ListRoutes()
		if err != nil {
			fmt.Println("Erro ao listar rotas:", err)
			return
		}

		for name, route := range routes {
			fmt.Printf("%s -> %d\n", name, route.Port)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
